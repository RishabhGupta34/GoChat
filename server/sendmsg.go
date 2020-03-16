package server

// Contains all the functions for sending messages to channel, user and through API

import (
	"fmt"
	"net/http"
	"time"
)

// SendMsgAPI is For sending a message to a channel using HTTP REST API
func SendMsgAPI(ServerVars *Variables, d DataAPI, w *http.ResponseWriter, t time.Time) {
	ServerVars.ChannelListMu.Lock()
	defer ServerVars.ChannelListMu.Unlock()

	_, exists := ServerVars.ChannelList[d.Channel] //Checking if the channel exists or not
	if !exists {
		_, err := (*w).Write([]byte(fmt.Sprintf("No channel exists with the name %s\n", d.Channel)))
		CheckError(ServerVars, err)
		return
	}

	if ServerVars.ChannelList[d.Channel].Access == "private" && d.Key != ServerVars.ChannelList[d.Channel].Key {
		_, err := (*w).Write([]byte(fmt.Sprintf("Key is not valid for the private channel %s\n", d.Channel)))
		CheckError(ServerVars, err)
		return
	}
	ServerVars.ChannelList[d.Channel].MessageHistChannel <- &Data{"API", d.Msg, t}

	// WriteMsgHistory(ServerVars, d.Channel, &Data{"API", d.Msg, t})

	ts := t.Format("2006-01-02 15:04:05")
	ServerVars.ActivityLogsChannel <- fmt.Sprintf(ServerVars.LogActivityFormat, "Send Message(Channel)", "", d.Channel, "API", ts)
	ServerVars.MessageLogsChannel <- fmt.Sprintf(ServerVars.LogMessageFormat, "Channel", d.Channel, "API", d.Msg, ts)

	for k := range ServerVars.ChannelList[d.Channel].UserList { //Iterating over members of the channel
		txt := fmt.Sprintf("\nChannel: %s\nMessage: %s\nSent by: API\nTime: %s\n\n", Colors["Cyan"](d.Channel), d.Msg, Colors["Yellow"](ts))
		_, err := (*ServerVars.ChannelList[d.Channel].UserList[k]).Write([]byte(txt)) // Sending msg to user "k"
		CheckError(ServerVars, err)
	}

	_, err := (*w).Write([]byte("Message sent\n"))
	CheckError(ServerVars, err)
}

// SendMsgC is For sending a message to a channel
func (c *Client) SendMsgC(ServerVars *Variables, ch string, d *Data) {
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
	ServerVars.ChannelList[ch].MessageHistChannel <- d
	// WriteMsgHistory(ServerVars, ch, d)
	t := d.Time.Format("2006-01-02 15:04:05")
	ServerVars.ActivityLogsChannel <- fmt.Sprintf(ServerVars.LogActivityFormat, "Send Message(Channel)", "", ch, c.Username, t) // Logging activity
	ServerVars.MessageLogsChannel <- fmt.Sprintf(ServerVars.LogMessageFormat, "Channel", ch, d.Username, d.Msg, t)              // Logging msg to log file

	ServerVars.BlockListMu.Lock()
	defer ServerVars.BlockListMu.Unlock()

	for k := range ServerVars.ChannelList[ch].UserList { // Iterating over members of the channel
		if k == d.Username {
			ConnWrite(ServerVars, c, "%s \"%s\" %s\n", Colors["Green"]("Message"), d.Msg, Colors["Green"]("sent"))
			continue
		}
		_, exists = ServerVars.BlockList[k][d.Username] // Checking if the user "k" has blocked the client(message sender) or not
		if exists {
			continue
		}
		txt := fmt.Sprintf("\nChannel: %s\nMessage: %s\nSent by: %s\nTime: %s\n\n", Colors["Cyan"](ch), d.Msg, Colors["Blue"](d.Username), Colors["Yellow"](t))
		_, err := (*ServerVars.ChannelList[ch].UserList[k]).Write([]byte(txt)) // Sending msg to user "k"
		CheckError(ServerVars, err)
	}
}

// SendMsgU is For sending a message to a specific user
func (c *Client) SendMsgU(ServerVars *Variables, u string, d *Data) {
	ServerVars.UsernamesMu.Lock()
	defer ServerVars.UsernamesMu.Unlock()

	_, exists := ServerVars.Usernames[u] // Checking if the user "u" exists or not
	if !exists {
		ConnWrite(ServerVars, c, Colors["Red"]("User doesn't exist\n"))
		return
	}

	ConnWrite(ServerVars, c, "%s \"%s\" %s\n", Colors["Green"]("Message"), d.Msg, Colors["Green"]("sent"))
	ServerVars.BlockListMu.Lock()
	defer ServerVars.BlockListMu.Unlock()

	_, exists = ServerVars.BlockList[u][d.Username] // Checking if the user "u" has blocked the client(message sender) or not
	if exists {
		return
	}

	t := d.Time.Format("2006-01-02 15:04:05")
	ServerVars.ActivityLogsChannel <- fmt.Sprintf(ServerVars.LogActivityFormat, "Send Message(User)", "", u, c.Username, t) // Logging activity
	ServerVars.MessageLogsChannel <- fmt.Sprintf(ServerVars.LogMessageFormat, "User", u, d.Username, d.Msg, t)              // Logging msg to log file

	txt := fmt.Sprintf("\nMessage: %s\nSent by: %s\nTime: %s\n\n", d.Msg, Colors["Blue"](d.Username), Colors["Yellow"](t))
	_, err := (*ServerVars.Usernames[u]).Write([]byte(txt)) // Sending msg to user "u"
	CheckError(ServerVars, err)
}
