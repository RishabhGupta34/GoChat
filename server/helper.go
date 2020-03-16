package server

// Contains all the helper functions

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

// ClearData is For clearing all the previous logs and message histories
func ClearData(ServerVars *Variables) {
	f := []string{ServerVars.Conf.Logs, ServerVars.Conf.ChDataFolder}
	for _, v := range f {
		err := os.RemoveAll(v)
		if err != nil && !os.IsNotExist(err) {
			log.Fatal(err) // Logging fatal errors
		}
	}
}

// Logging is for checking the all the log channels and pushing them to respective log files
func Logging(ServerVars *Variables) {
	for {
		select {
		case e, more := <-ServerVars.ErrorLogsChannel: // for checking the ErrorLogsChannel and pushing errors to error log file
			if !more {
				ServerVars.ErrorLogsChannel = nil
			} else {
				ServerVars.LogError.Println(e)
			}
		case m, more := <-ServerVars.MessageLogsChannel: // For checking the MessageLogsChannel and pushing messages to message log file
			if !more {
				ServerVars.MessageLogsChannel = nil
			} else {
				ServerVars.LogMessage.Println(m)
			}
		case a, more := <-ServerVars.ActivityLogsChannel: // For checking the ActivityLogsChannel and pushing activity to activity log file
			if !more {
				ServerVars.ActivityLogsChannel = nil
			} else {
				ServerVars.LogActivity.Println(a)
			}
		}
		if ServerVars.ErrorLogsChannel == nil && ServerVars.MessageLogsChannel == nil && ServerVars.ActivityLogsChannel == nil {
			break // If all 3 channels are closed break the infinite loop
		}
	}
}

// CheckError is for checking if the error is not nil and pushing them to error log channel
func CheckError(ServerVars *Variables, err error) {
	if err != nil {
		ServerVars.ErrorLogsChannel <- err
	}
}

// LoggerF is for logging fatal errors if error is not nil
func LoggerF(ServerVars *Variables, err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// ConnWrite is func for writing to the client's connection
func ConnWrite(ServerVars *Variables, c *Client, format string, a ...interface{}) {
	_, err := (*c.Conn).Write([]byte(fmt.Sprintf(format, a...)))
	CheckError(ServerVars, err)
}

// InterruptHandler is for catching interrupts and then closing all the channels
func InterruptHandler(ServerVars *Variables, opt string) { // opt is for diff. b/w main and test
	select {
	case <-ServerVars.SignalChannel:
		close(ServerVars.ErrorLogsChannel)
		close(ServerVars.MessageLogsChannel)
		close(ServerVars.ActivityLogsChannel)
		for _, v := range ServerVars.ChannelList {
			close(v.MessageHistChannel)
		}
		if opt == "main" {
			fmt.Println("\nSHUTTING DOWN GoChat!")
			os.Exit(0)
		}
	}
}

// WriteMsgHistory writes the messages that are sent to each channel separately
func (c *Channels) WriteMsgHistory(ServerVars *Variables) {
	for {
		d, more := <-c.MessageHistChannel
		if !more {
			break // If message hist channels is closed break the infinite loop
		}
		c.Mu.Lock()
		bt, err := ioutil.ReadFile(fmt.Sprintf("%s/%s.json", ServerVars.Conf.ChDataFolder, c.Name)) // Reading the channel message history for updating it
		CheckError(ServerVars, err)

		var Dhist []Data
		err = json.Unmarshal(bt, &Dhist)
		CheckError(ServerVars, err)

		Dhist = append(Dhist, *d) // Appending the new data to the history
		bt, err = json.Marshal(Dhist)
		CheckError(ServerVars, err)

		err = ioutil.WriteFile(fmt.Sprintf("%s/%s.json", ServerVars.Conf.ChDataFolder, c.Name), bt, 0766)
		CheckError(ServerVars, err)
		
		c.Mu.Unlock()
	}
}
