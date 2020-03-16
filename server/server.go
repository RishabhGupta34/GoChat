package server

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"time"
)

// Initialise is For initialising all the global variables
func Initialise(ServerVars *Variables, cp string, opt string) { // opt is for diff. b/w main and test
	ServerVars.LogMessageFormat = "To%s:- %s\n\t\t\t\tSender:- %s\n\t\t\t\tMessage:- %s \n\t\t\t\tTime:- %s\n"
	ServerVars.LogActivityFormat = "Activity: %s%s\n\t\t\t\tName: %s\n\t\t\t\tBy: %s\n\t\t\t\tTime: %s\n"
	ServerVars.BlockList = make(map[string]map[string]bool)
	ServerVars.ChannelList = make(map[string]*Channels)
	ServerVars.Usernames = make(map[string]*net.Conn)
	ServerVars.BlockListMu = sync.Mutex{}
	ServerVars.ChannelListMu = sync.Mutex{}
	ServerVars.UsernamesMu = sync.Mutex{}
	ServerVars.ErrorLogsChannel = make(chan error, 100)
	ServerVars.MessageLogsChannel = make(chan string, 100)
	ServerVars.ActivityLogsChannel = make(chan string, 100)
	ServerVars.SignalChannel = make(chan os.Signal)
	signal.Notify(ServerVars.SignalChannel, os.Interrupt)
	go InterruptHandler(ServerVars, opt)

	// if len(os.Args) != la { // For checking command line argument length and it is diff. for main and test
	// 	log.Fatal("Please enter the full path for the server's config file")
	// }

	go Logging(ServerVars)

	bt, err := ioutil.ReadFile(cp) // Reading the config file
	LoggerF(ServerVars, err)

	err = json.Unmarshal(bt, &ServerVars.Conf)
	LoggerF(ServerVars, err)

	ClearData(ServerVars)                      // Clearing all previous data
	err = os.Mkdir(ServerVars.Conf.Logs, 0777) // Creating Logs folder
	if err != nil && os.IsNotExist(err) {
		log.Fatal(err) // Logging fatal errors
	}

	err = os.Mkdir(ServerVars.Conf.ChDataFolder, 0777) // Creating Channel folder
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
