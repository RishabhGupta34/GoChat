package server

// Contains the functions for creating a new client and deleting a client

import (
	"net"
	"strings"
	"time"
)

// DeleteUser is Function for deleting a user i.e when a client stops using the chat server
func (c *Client) DeleteUser(ServerVars *Variables, t time.Time) {
	ConnWrite(ServerVars, c, "Thanks for using GoChat %s \n", Colors["Blue"](c.Username)) // Final greeting message
	ServerVars.UsernamesMu.Lock()
	defer ServerVars.UsernamesMu.Unlock()
	ServerVars.BlockListMu.Lock()
	defer ServerVars.BlockListMu.Unlock()
	ServerVars.ChannelListMu.Lock()
	defer ServerVars.ChannelListMu.Unlock()
	delete(ServerVars.Usernames, c.Username) // Deleting the username of the client from the list of used usernames
	delete(ServerVars.BlockList, c.Username) // Deleting the blocklist of the client
	for _, v := range c.Channels {
		delete(ServerVars.ChannelList[v].UserList, c.Username) //Deleting the client from all the subscribed channels
	}
	for k := range ServerVars.BlockList {
		_, exists := ServerVars.BlockList[k][c.Username]
		if exists { //Deleting the client from the blocklist of all users so that if a new user with same username joins he doesn't get blocked
			delete(ServerVars.BlockList[k], c.Username)
		}
	}
	ServerVars.LogActivity.Printf(ServerVars.LogActivityFormat, "Delete User", "", c.Username, c.Username, t.Format("2006-01-02 15:04:05")) // Logging activity
}

// NewClient is For making a new client
func NewClient(ServerVars *Variables, conn *net.Conn, t time.Time) Client {
	var u string // Saves the username of the client
	ServerVars.UsernamesMu.Lock()
	defer ServerVars.UsernamesMu.Unlock()
	for { // Runs until the client enters a valid username
		_, err := (*conn).Write([]byte("Enter name: "))
		LoggerE(ServerVars, err)
		buff := make([]byte, 100)
		n, err := (*conn).Read(buff)
		if err != nil {
			ServerVars.LogError.Println(err) // Logging errors
			continue
		}
		u = string(buff[:n])
		u = strings.Trim(u, "\r\n")
		_, exists := ServerVars.Usernames[u] // Checking if the username already exists or not
		if exists {
			_, err = (*conn).Write([]byte(Colors["Red"]("Username already exists!!\n")))
			LoggerE(ServerVars, err)
			continue
		} else {
			break
		}
	}
	ServerVars.LogActivity.Printf(ServerVars.LogActivityFormat, "New User", "", u, u, t.Format("2006-01-02 15:04:05")) // Logging activity
	ServerVars.Usernames[u] = conn
	c := Client{
		Username: u,
		Channels: make([]string, 0),
		Conn:     conn,
	}
	ConnWrite(ServerVars, &c, "%s %s\n\n", "Hey", Colors["Blue"](c.Username))
	c.JoinChannel(ServerVars, "all", "nil", time.Now()) // Joining the channel ALL
	ServerVars.BlockListMu.Lock()
	defer ServerVars.BlockListMu.Unlock()
	ServerVars.BlockList[c.Username] = make(map[string]bool) // Creating a blocklist for the client
	return c
}
