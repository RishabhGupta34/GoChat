package server

// Contains the function for printing the available options to the client

// PrintOptions is To print the available options
func (c *Client) PrintOptions(ServerVars *Variables) {
	ConnWrite(ServerVars, c, "\n\\Options : Show all the available options\n")
	ConnWrite(ServerVars, c, "\\Channels : Show all the subscribed channels\n")
	ConnWrite(ServerVars, c, "\\Info <channel_name> : Show the channel's information\n")
	ConnWrite(ServerVars, c, "\\New <channel_name> <access> <key>: Create a new channel with public/private access\n")
	ConnWrite(ServerVars, c, "\\Join <channel_name> <key>: Subscribe to <channel_name>\n")
	ConnWrite(ServerVars, c, "\\Leave <channel_name> : Un-subscribe from <channel_name>\n")
	ConnWrite(ServerVars, c, "\\SendC <channel_name> <Message> : Send <Message> to <channel_name>\n")
	ConnWrite(ServerVars, c, "\\SendU <user_name> <Message> : Send <Message> to <user_name>\n")
	ConnWrite(ServerVars, c, "\\Block <username> : Block <username>\n")
	ConnWrite(ServerVars, c, "\\Unblock <username> : Unblock <username>\n")
	ConnWrite(ServerVars, c, "\\Blocklist : List of blocked users\n")
	ConnWrite(ServerVars, c, "\\Close : Close the connection\n\n")
	ConnWrite(ServerVars, c, "Use channel %s to send message to all the users\n", Colors["Cyan"]("all"))
	ConnWrite(ServerVars, c, "Key is a unique string that will be used to join a channel. Enter nil if channel is public\n")
	ConnWrite(ServerVars, c, Colors["Red"]("Channel name is not case sensitive\n\n\n"))
}
