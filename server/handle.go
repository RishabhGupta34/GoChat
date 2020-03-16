package server

// Contains the function that handles the client's connection

import (
	"net"
	"strings"
	"time"
)

// HandleConn is For handling a connection
func HandleConn(ServerVars *Variables, conn *net.Conn) {
	c := NewClient(ServerVars, conn, time.Now()) // Returns Client struct
	c.PrintOptions(ServerVars)                   // Printing the available options

	flag := true // Saves if the connection is open or closed by the client
	for flag {
		buff := make([]byte, 1000)
		n, err := (*c.Conn).Read(buff) // Reading input from the user
		if err != nil {
			ServerVars.ErrorLogsChannel <- err
			continue
		}

		currTime := time.Now()
		comm := string(buff[:n])          // Converting bytes to string
		comm = strings.Trim(comm, "\r\n") // Trimming new line characters
		comm = strings.TrimSpace(comm)    // Trimming extra space
		splComm := strings.Split(comm, " ")

		switch strings.ToLower(splComm[0]) {
		case "\\options": // Printing the available options
			c.PrintOptions(ServerVars)

		case "\\channels": // Printing the subscribed channels of the client
			for _, v := range c.Channels {
				ConnWrite(ServerVars, &c, "%s\n", Colors["Cyan"](v))
			}

		case "\\info": // Printing information about a channel
			if len(splComm) != 2 {
				ConnWrite(ServerVars, &c, Colors["Red"]("Wrong option!\n"))
			} else {
				c.Info(ServerVars, strings.ToLower(splComm[1]))
			}

		case "\\new": // Creating a new channel
			if len(splComm) != 4 {
				ConnWrite(ServerVars, &c, Colors["Red"]("Wrong option!\n"))
			} else {
				NewChannel(ServerVars, strings.ToLower(splComm[1]), splComm[2], splComm[3], &c, currTime)
			}

		case "\\join": // Joining a channel
			if len(splComm) != 3 {
				ConnWrite(ServerVars, &c, Colors["Red"]("Wrong option!\n"))
			} else {
				c.JoinChannel(ServerVars, strings.ToLower(splComm[1]), splComm[2], currTime)
			}

		case "\\leave": // Leaving a channel
			if len(splComm) != 2 {
				ConnWrite(ServerVars, &c, Colors["Red"]("Wrong option!\n"))
			} else {
				c.LeaveChannel(ServerVars, strings.ToLower(splComm[1]), currTime)
			}

		case "\\sendc": // Sending message to a channel
			if len(splComm) < 3 {
				ConnWrite(ServerVars, &c, Colors["Red"]("Wrong option!\n"))
			} else {
				d := Data{c.Username, strings.Join(splComm[2:], " "), currTime}
				c.SendMsgC(ServerVars, strings.ToLower(splComm[1]), &d)
			}

		case "\\sendu": // Sending message to a user
			if len(splComm) < 3 {
				ConnWrite(ServerVars, &c, Colors["Red"]("Wrong option!\n"))
			} else {
				d := Data{c.Username, strings.Join(splComm[2:], " "), currTime}
				c.SendMsgU(ServerVars, splComm[1], &d)
			}

		case "\\block": // Blocking a user
			if len(splComm) != 2 {
				ConnWrite(ServerVars, &c, Colors["Red"]("Wrong option!\n"))
			} else {
				c.BlockUser(ServerVars, splComm[1], currTime)
			}

		case "\\unblock": // Unblocking a user
			if len(splComm) != 2 {
				ConnWrite(ServerVars, &c, Colors["Red"]("Wrong option!\n"))
			} else {
				c.UnblockUser(ServerVars, splComm[1], currTime)
			}

		case "\\blocklist": // Printing the list of blocked users
			if len(ServerVars.BlockList[c.Username]) == 0 {
				ConnWrite(ServerVars, &c, Colors["Red"]("You have not blocked anyone\n"))
			}
			for k := range ServerVars.BlockList[c.Username] {
				ConnWrite(ServerVars, &c, "%s\n", Colors["Cyan"](k))
			}

		case "\\close": // Closing a connection
			c.DeleteUser(ServerVars, currTime) // Delete the client and all its information
			err := (*c.Conn).Close()
			CheckError(ServerVars, err)
			flag = false // For breaking the loop

		case "": // Do nothing on empty string
		default:
			ConnWrite(ServerVars, &c, Colors["Red"]("Wrong option!\n"))
		}
	}
}
