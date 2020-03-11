package server

// Contains all the variables and structures for the server

import (
	"log"
	"net"
	"sync"
	"time"
)

// Config is This struct saves the configuration of the chat server
type Config struct {
	Host         string `json:"host"`         // Saves hostname
	APIPort      string `json:"apiport"`      // Saves port number for api calls
	TelPort      string `json:"telport"`      // Saves port number for telnet
	ErrorLog     string `json:"error-log"`    // Filepath to save the error logs
	MsgLog       string `json:"msg-log"`      // Filepath to save the message logs
	ActLog       string `json:"act-log"`      // Filepath to save the activity logs
	Logs         string `json:"logs"`         // Filepath to save the activity logs
	ChDataFolder string `json:"channel-data"` // Folder path to save the channel message history
}

// Data is Saves the message data that is sent by any client
type Data struct {
	Username string    // Username of the client that sent the message
	Msg      string    // The actual msg
	Time     time.Time // Time at which the msg was sent
}

// Client is Saves information about a client
type Client struct {
	Username string    // Unique username of a client
	Channels []string  // Channel subscriptions of a client
	Conn     *net.Conn // Saves the network connection of this client
}

// DataAPI is Saves information that is sent through HTTP REST POST API calls
type DataAPI struct {
	Channel string // Name of the channel
	Key     string // If channel is private this unique key is used to send msg to the channel
	Msg     string // The actual msg
}

// Channels is Saves the channel information
type Channels struct {
	Name      string               // Name of the channel
	Creator   string               // Creator(Username) of the channel
	CreatedAt time.Time            // Time of creation of the channel
	UserList  map[string]*net.Conn // Mapping of usernames(that are in this channel) to their network connection
	Access    string               // Type of channel, can be public or private
	Key       string               // If channel is private this unique key is used to join the channel
}

// Mutexes contains mutexes for the shared variables
type Mutexes struct {
	BlockList   *sync.Mutex // Mutex for Blocklist
	ChannelList *sync.Mutex // Mutex for ChannelList
	Usernames   *sync.Mutex // Mutex for Usernames
}

// Variables contains all the variables for the server
type Variables struct {
	LogMessageFormat  string                     // LogMessageFormat is Default message for logging messages
	LogActivityFormat string                     // LogActivityFormat is Default message for logging activity
	LogError          *log.Logger                // LogError is Logger for logging errors
	LogMessage        *log.Logger                // LogMessage is Logger for logging messages
	LogActivity       *log.Logger                // LogActivity is Logger for logging all activity(channel creation,new user,joining channel etc) on the server
	Conf              Config                     // Conf is Saves the configuration of the chat server
	BlockList         map[string]map[string]bool /* BlockList is Contains the blocked users list of each user
	 key -> specifying the username i.e whose blocklist is it.
	value -> map containing all the blocked users of the "key" */

	ChannelList map[string]Channels /* ChannelList is Contains list of all the channels
	Mapping from channel-name -> Channels struct(contain info about channel) */

	Usernames     map[string]*net.Conn // Usernames is Saves all the used usernames
	BlockListMu   *sync.Mutex          // Mutex for Blocklist
	ChannelListMu *sync.Mutex          // Mutex for ChannelList
	UsernamesMu   *sync.Mutex          // Mutex for Usernames
	// Colors            map[string]func(...interface{}) string
}

// Colors is map containing different functions for coloring the text shown to the client
var Colors map[string]func(...interface{}) string
