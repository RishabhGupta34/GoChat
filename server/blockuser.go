package server

// Contains functions that handles the blocking and unblocking a user

import (
	"time"
)

// BlockUser is For blocking a specific user
func (c *Client) BlockUser(ServerVars *Variables, u string, t time.Time) {
	if c.Username == u { // Checking if the client is trying to block himself
		ConnWrite(ServerVars, c, "%s\n", Colors["Red"]("You cannot block yourself"))
		return
	}
	ServerVars.UsernamesMu.Lock()
	defer ServerVars.UsernamesMu.Unlock()
	_, exists := ServerVars.Usernames[u] // Checking if the user "u" exists or not
	if !exists {
		ConnWrite(ServerVars, c, "%s\n", Colors["Red"]("User doesn't exist"))
		return
	}
	ServerVars.BlockListMu.Lock()
	defer ServerVars.BlockListMu.Unlock()
	_, exists = ServerVars.BlockList[c.Username][u] // Checking if the user "u" already exists in the client's blocklist
	if exists {
		ConnWrite(ServerVars, c, "%s %s %s\n", Colors["Red"]("User"), Colors["Blue"](u), Colors["Red"]("already belongs in your blocklist"))
		return
	}
	ServerVars.LogActivity.Printf(ServerVars.LogActivityFormat, "Block User", "", u, c.Username, t.Format("2006-01-02 15:04:05")) // Logging activity
	ConnWrite(ServerVars, c, "%s %s\n", Colors["Green"]("Blocked user"), Colors["Blue"](u))
	ServerVars.BlockList[c.Username][u] = true // Adding the user "u" to the client's blocklist
}

// UnblockUser is For un-blocking a specific user
func (c *Client) UnblockUser(ServerVars *Variables, u string, t time.Time) {
	ServerVars.UsernamesMu.Lock()
	defer ServerVars.UsernamesMu.Unlock()
	_, exists := ServerVars.Usernames[u] // Checking if the user "u" exists or not
	if !exists {
		ConnWrite(ServerVars, c, "%s\n", Colors["Red"]("User doesn't exist"))
		return
	}
	ServerVars.BlockListMu.Lock()
	defer ServerVars.BlockListMu.Unlock()
	_, exists = ServerVars.BlockList[c.Username][u] // Checking if the user "u" doesn't exist in the client's blocklist
	if !exists {
		ConnWrite(ServerVars, c, "%s %s %s\n", Colors["Red"]("User"), Colors["Blue"](u), Colors["Red"]("doesn't belong in your blocklist"))
		return
	}
	ServerVars.LogActivity.Printf(ServerVars.LogActivityFormat, "Un-Block User", "", u, c.Username, t.Format("2006-01-02 15:04:05")) // Logging activity
	ConnWrite(ServerVars, c, "%s %s\n", Colors["Green"]("Unblocked user"), Colors["Blue"](u))
	delete(ServerVars.BlockList[c.Username], u) // Removing the user "u" from the client's blocklist
}
