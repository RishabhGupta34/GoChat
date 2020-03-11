package server

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net"
	"os"
	"sync"
	"time"

	"github.com/fatih/color"
)

// Initialise is For initialising all the global variables
func Initialise(ServerVars *Variables, la int, i int) {
	ServerVars.LogMessageFormat = "To%s:- %s\n\t\t\t\tSender:- %s\n\t\t\t\tMessage:- %s \n\t\t\t\tTime:- %s\n\n"
	ServerVars.LogActivityFormat = "Activity: %s%s\n\t\t\t\tName: %s\n\t\t\t\tBy: %s\n\t\t\t\tTime: %s\n\n"
	ServerVars.BlockList = make(map[string]map[string]bool)
	ServerVars.ChannelList = make(map[string]Channels)
	ServerVars.Usernames = make(map[string]*net.Conn)
	ServerVars.BlockListMu = &sync.Mutex{}
	ServerVars.ChannelListMu = &sync.Mutex{}
	ServerVars.UsernamesMu = &sync.Mutex{}
	Colors = map[string]func(...interface{}) string{
		"Yellow": color.New(color.FgYellow, color.Bold).SprintFunc(),
		"Blue":   color.New(color.FgBlue, color.Bold).SprintFunc(),
		"Red":    color.New(color.FgRed, color.Bold).SprintFunc(),
		"Green":  color.New(color.FgGreen, color.Bold).SprintFunc(),
		"Cyan":   color.New(color.FgCyan, color.Bold).SprintFunc(),
	}
	if len(os.Args) != la {
		log.Fatal("Please enter the full path for the server's config file")
	}
	bt, err := ioutil.ReadFile(os.Args[i]) // Reading the config file
	LoggerF(ServerVars, err)
	err = json.Unmarshal(bt, &ServerVars.Conf)
	LoggerF(ServerVars, err)
	ClearData(ServerVars) // Clearing all previous data
	err = os.Mkdir(ServerVars.Conf.Logs, 0777)
	if err != nil && os.IsNotExist(err) {
		log.Fatal(err) // Logging fatal errors
	}
	err = os.Mkdir(ServerVars.Conf.ChDataFolder, 0777)
	if err != nil && os.IsNotExist(err) {
		log.Fatal(err) // Logging fatal errors
	}
	errlogf, err := os.OpenFile(ServerVars.Conf.ErrorLog, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	LoggerF(ServerVars, err)
	msglogf, err := os.OpenFile(ServerVars.Conf.MsgLog, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	LoggerF(ServerVars, err)
	actlogf, err := os.OpenFile(ServerVars.Conf.ActLog, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	LoggerF(ServerVars, err)
	ServerVars.LogError = log.New(errlogf, "", log.LstdFlags)     // Logger for logging errors
	ServerVars.LogMessage = log.New(msglogf, "", log.Lshortfile)  // Logger for logging messages
	ServerVars.LogActivity = log.New(actlogf, "", log.Lshortfile) // Logger for logging all the activity on the server

	NewChannel(ServerVars, "all", "public", "nil", nil, time.Now()) // Creating a new channel "ALL" for sending messages to all connected clients
}
