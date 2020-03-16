package server

// Contains the function for handling the read message history through API calls

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// ReadMsgAPI is For reading message history of a channel using HTTP REST API
func ReadMsgAPI(ServerVars *Variables, c string, k string, w *http.ResponseWriter, t time.Time) {
	ServerVars.ChannelListMu.Lock()
	defer ServerVars.ChannelListMu.Unlock()

	_, exists := ServerVars.ChannelList[c] //Checking if the channel exists or not
	if !exists {
		_, err := (*w).Write([]byte(fmt.Sprintf("No channel exists with the name %s\n", c)))
		CheckError(ServerVars, err)
		return
	}

	if ServerVars.ChannelList[c].Access == "private" && k != ServerVars.ChannelList[c].Key { // If the channnel is private, verify if key is correct or not
		_, err := (*w).Write([]byte(fmt.Sprintf("Key is not valid for the private channel %s\n", c)))
		CheckError(ServerVars, err)
		return
	}

	ServerVars.ChannelList[c].Mu.Lock()
	bt, err := ioutil.ReadFile(fmt.Sprintf("%s/%s.json", ServerVars.Conf.ChDataFolder, c)) // Reading json file containing the history
	CheckError(ServerVars, err)
	ServerVars.ChannelList[c].Mu.Unlock()

	var Dhist []Data
	err = json.Unmarshal(bt, &Dhist) // Converting json encoding to slice of type Data
	CheckError(ServerVars, err)

	ts := t.Format("2006-01-02 15:04:05")
	ServerVars.ActivityLogsChannel <- fmt.Sprintf(ServerVars.LogActivityFormat, "Read Message History", "", c, "API", ts) // Log Activity

	for k := range Dhist { //Iterating over messages
		ts = Dhist[k].Time.Format("2006-01-02 15:04:05")
		txt := fmt.Sprintf("\nMessage: %s\nSent by: %s\nTime: %s\n\n", Dhist[k].Msg, Colors["Blue"](Dhist[k].Username), Colors["Yellow"](ts))
		_, err := (*w).Write([]byte(txt)) // Sending msg to user "k"
		CheckError(ServerVars, err)
	}
}
