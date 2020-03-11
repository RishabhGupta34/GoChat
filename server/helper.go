package server

// Contains all the helper functions

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

//ClearData is For clearing all the previous logs and message histories
func ClearData(ServerVars *Variables) {
	f := []string{ServerVars.Conf.Logs, ServerVars.Conf.ChDataFolder}
	for _, v := range f {
		err := os.RemoveAll(v)
		if err != nil && !os.IsNotExist(err) {
			log.Fatal(err) // Logging fatal errors
		}
	}
}

//LoggerE is for logging errors if error is not nil
func LoggerE(ServerVars *Variables, err error) {
	if err != nil {
		ServerVars.LogError.Println(err)
	}
}

//LoggerF is for logging fatal errors if error is not nil
func LoggerF(ServerVars *Variables, err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// ConnWrite is func for writing to the client's connection
func ConnWrite(ServerVars *Variables, c *Client, format string, a ...interface{}) {
	_, err := (*c.Conn).Write([]byte(fmt.Sprintf(format, a...)))
	LoggerE(ServerVars, err)
}

// WriteMsgHistory writes the messages that are sent to each channel separately
func WriteMsgHistory(ServerVars *Variables, ch string, d *Data) {
	bt, err := ioutil.ReadFile(fmt.Sprintf("%s/%s.json", ServerVars.Conf.ChDataFolder, ch)) // Reading the channel message history for updating it
	LoggerE(ServerVars, err)
	var Dhist []Data
	err = json.Unmarshal(bt, &Dhist)
	LoggerE(ServerVars, err)
	Dhist = append(Dhist, *d) // Appending the new data to the history
	bt, err = json.Marshal(Dhist)
	LoggerE(ServerVars, err)
	err = ioutil.WriteFile(fmt.Sprintf("%s/%s.json", ServerVars.Conf.ChDataFolder, ch), bt, 0766)
	LoggerE(ServerVars, err)
}
