package server

// Contains functions for channels: new channel, join channel, leave channel, printing info of a channel
import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"strings"
	"time"
)

// NewChannel is Creating a new channel
func NewChannel(ServerVars *Variables, ch string, a string, k string, c *Client, t time.Time) {
	ServerVars.ChannelListMu.Lock()
	defer ServerVars.ChannelListMu.Unlock()
	if c == nil { // If the admin of the server wants to create a channel
		ServerVars.ChannelList[ch] = Channels{ch, "Admin", t, make(map[string]*net.Conn), a, k}
		ServerVars.LogActivity.Printf(ServerVars.LogActivityFormat, "Channel creation ", strings.ToUpper(a), ch, "Admin", t.Format("2006-01-02 15:04:05"))
		return
	}
	if (a != "public") && (a != "private") { // Checking if the access type specified by the client is correct or not
		ConnWrite(ServerVars, c, Colors["Red"]("Invalid access type\n"))
		return
	}
	_, exists := ServerVars.ChannelList[ch]
	if exists { // Checking if the channel already exists
		ConnWrite(ServerVars, c, "%s %s\n", Colors["Red"]("Channel exists with the name:"), ch)
		return
	}
	var d []Data // Type of data stored for retrieving message history
	bt, err := json.Marshal(d)
	LoggerE(ServerVars, err)
	err = ioutil.WriteFile(fmt.Sprintf("%s/%s.json", ServerVars.Conf.ChDataFolder, ch), bt, 0766) // Creating a json file for storing all the message history of a channel
	LoggerE(ServerVars, err)
	if a == "public" {
		k = "nil" // Setting the key to be nil if the channel is public
	}
	c.Channels = append(c.Channels, ch)                                                        // Adding the channel to the client's list of subscribed channels
	ServerVars.ChannelList[ch] = Channels{ch, c.Username, t, make(map[string]*net.Conn), a, k} // Adding the channel to the list of all channels
	ConnWrite(ServerVars, c, "%s %s %s %s %s\n", Colors["Green"]("New "), Colors["Green"](a), Colors["Green"](" channel"), Colors["Cyan"](ch), Colors["Green"]("created"))
	ServerVars.ChannelList[ch].UserList[c.Username] = c.Conn // Adding the client to the members list of the channel
	ServerVars.LogActivity.Printf(ServerVars.LogActivityFormat, "Channel creation ", strings.ToUpper(a), ch, c.Username, t.Format("2006-01-02 15:04:05"))
}

// JoinChannel is Function for adding a user to a new channel
func (c *Client) JoinChannel(ServerVars *Variables, ch string, k string, t time.Time) {
	ServerVars.ChannelListMu.Lock()
	defer ServerVars.ChannelListMu.Unlock()
	_, exists := ServerVars.ChannelList[ch]
	if !exists { // For checking if the specified channel exists or not
		ConnWrite(ServerVars, c, "%s %s\n", Colors["Red"]("No channel exists with the name:"), Colors["Cyan"](ch))
		return
	}
	for _, v := range c.Channels {
		if v == ch { // For checking if the client is already a member of the specified channel
			ConnWrite(ServerVars, c, "%s %s\n", Colors["Red"]("You have already joined the channel"), Colors["Cyan"](ch))
			return
		}
	}
	if ServerVars.ChannelList[ch].Access == "private" && ServerVars.ChannelList[ch].Key != k { // If the channnel is private, verify if key is correct or not
		ConnWrite(ServerVars, c, Colors["Red"]("Please enter valid key for joining a private channel\n"))
		return
	}
	c.Channels = append(c.Channels, ch) // Adding the channel to the client's list of subscribed channels
	ConnWrite(ServerVars, c, "%s %s\n", Colors["Green"]("You were added to channel"), Colors["Cyan"](ch))
	ServerVars.ChannelList[ch].UserList[c.Username] = c.Conn // Adding the client to the members list of the channel
	ServerVars.LogActivity.Printf(ServerVars.LogActivityFormat, "Join Channel", "", ch, c.Username, t.Format("2006-01-02 15:04:05"))
}

// LeaveChannel is For leaving a channel
func (c *Client) LeaveChannel(ServerVars *Variables, ch string, t time.Time) {
	ServerVars.ChannelListMu.Lock()
	defer ServerVars.ChannelListMu.Unlock()
	if ch == "all" { // Client cannot leave the channel all
		ConnWrite(ServerVars, c, "%s %s\n", Colors["Red"]("You cannot leave the channel"), Colors["Cyan"](ch))
		return
	}
	_, exists := ServerVars.ChannelList[ch] // Checking if the channel "ch" exists or not
	if !exists {
		ConnWrite(ServerVars, c, "%s %s\n", Colors["Red"]("You are not a member of the channel"), Colors["Cyan"](ch))
		return
	}
	_, exists = ServerVars.ChannelList[ch].UserList[c.Username] // Checking if the client is a member of the channel "ch" or not
	if !exists {
		ConnWrite(ServerVars, c, "%s %s\n", Colors["Red"]("You are not a member of the channel"), Colors["Cyan"](ch))
		return
	}
	delete(ServerVars.ChannelList[ch].UserList, c.Username) // Removing the client from the userlist of the channel
	for i, v := range c.Channels {
		if v == ch {
			c.Channels = append(c.Channels[:i], c.Channels[i+1:]...) // Removing the channel from client's list of joined channels
			break
		}
	}
	ConnWrite(ServerVars, c, "%s %s\n", Colors["Green"]("You have left"), Colors["Cyan"](ch))
	ServerVars.LogActivity.Printf(ServerVars.LogActivityFormat, "Leave Channel", "", ch, c.Username, t.Format("2006-01-02 15:04:05")) // Log Activity
}

// Info is To show the information about a channel
func (c *Client) Info(ServerVars *Variables, ch string) {
	ServerVars.ChannelListMu.Lock()
	defer ServerVars.ChannelListMu.Unlock()
	_, exists := ServerVars.ChannelList[ch] //Checking if the channel "ch" exists or not
	if !exists {
		ConnWrite(ServerVars, c, "%s %s\n", Colors["Red"]("You have not joined the channel"), Colors["Cyan"](ch))
		return
	}
	_, exists = ServerVars.ChannelList[ch].UserList[c.Username] //Checking if the client is a member of the channel "ch" or not
	if !exists {
		ConnWrite(ServerVars, c, "%s %s\n", Colors["Red"]("You have not joined the channel"), Colors["Cyan"](ch))
		return
	}
	t := ServerVars.ChannelList[ch].CreatedAt.Format("2006-01-02 15:04:05")
	ServerVars.LogActivity.Printf(ServerVars.LogActivityFormat, "Print Channel Info", "", ch, c.Username, t) // Log Activity
	txt := "Channel name:- %s\nAccess:- %s\nKey:- %s\nCreated by:- %s\nCreated At:- %s\nMembers:-\n"
	a := ServerVars.ChannelList[ch].Access
	k := ServerVars.ChannelList[ch].Key
	cr := ServerVars.ChannelList[ch].Creator
	ConnWrite(ServerVars, c, txt, Colors["Cyan"](ch), a, k, Colors["Blue"](cr), Colors["Yellow"](t))
	for k := range ServerVars.ChannelList[ch].UserList { // Printing the members of the channel "ch"
		ConnWrite(ServerVars, c, "\t%s\n", Colors["Blue"](k))
	}
}
