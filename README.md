# GoChat
Chat server written in Golang. Multiple clients can connect via telnet. HTTP REST API to send messages and query for channel message history. 

## How to run the server:
Place this whole folder inside the GOPATH

	   go run server.go <config_file_full_path>

Change the config.json file to change the configuration (like host, port number, log file locations) of the server.


## Requirements:
Go module will install the packages required


## Features:
1.	There can be multiple connected clients at the same time. 
2.	All the connected clients at a time have unique usernames.
3.	Usernames can be reused (after the client with same username has closed the connection) 
4.	Client can send a direct message to another client provided that he has the correct username of the client.
5.	Any client can create a new channel provided a channel with same name doesn’t exist already.
6.	Both public and private Channels can be created

	    •	Private channels can only be joined by clients which have the unique key of the channel.

	    •	Public channels can be joined by anyone. Mention key as nil.
7.	A client can send message to all connected clients using the channel “all”.
8.	A client can block and unblock a specific client i.e. a client can choose to unsubscribe from another client’s messages.
9.	A client can view its own blocklist.
10.	A client can view the channel information:

	    •	Creator of the channel
	
	    •	Time the channel was created
	
	    •	Channel access(public/private)
	
	    •	Channel key 
	
	    •	Members of the channel
	
	This information can only be viewed if a client has joined the channel.
11.	Client can view the channels that he/she has joined.
12.	HTTP REST API for posting messages to a channel.
13.	HTTP REST API for reading message history of a channel (Requires key if the channel is private).
14.	A configuration file saves all the configuration (like host, port number, log file locations) of the chat server
15.	Error logging to log all the errors.
16.	Message logging to log all the messages (both to channel and direct messages)
17.	Activity logging to log activities on the chat server:

	    •	New User

	    •	Channel Creation
    
	    •	Join Channel
    
	    •	Leave Channel
    
	    •	Send Message to a channel
    
	    •	Send Direct Message
    
	    •	Block User
    
	    •	Unblock User
    
	    •	Read message history (API) 
    
 	    •	Delete User

Unit testing is also implemented with 88% coverage.

## Limitations:
1.	Use of local database to store log files and message history.
2.	Channel name is not case sensitive.
3.	No concept of channel admin
4.	If a client is typing a command (to perform some action) in the telnet server and another client sends him a message, he has to type the command again


## Examples:
1.	\options – To print the available options
2.	\channels – To print the list of client’s subscribed channels
3.	\info all – To print the information about the channel “all”
4.	\new tech w2ci93 private – To create a private channel “tech” with the key “w2ci93”
5.	\new gen nil public – To create a public channel “gen” with the key “nil”
6.	\join tech w2ci93 – To join the private channel “tech”. Key should be valid to join a private channel.
7.	\join gen nil - To join the public channel “gen”.
8.	\leave tech – To leave the channel “tech”
9.	\sendc gen Welcome to the channel – To send a message “Welcome to the channel” to channel “gen”
10.	\sendu rish Hey rish – To send a direct message “Hey rish” to the user “rish”
11.	\block rish – To block the user “rish”
12.	\unblock rish – To unblock the user “rish”
13.	\blocklist – To print the blocklist of the client
14.	\close – To close the connection 
15.	curl -X POST -d "{'channel': 'tech','key':'nil','msg':'Hello Hi'}" localhost:8080/send – To send a message to channel gen
16.	curl -X GET localhost:8080/history/tech/w2ci93 – To get message history of the channel tech. Key should be specified in the last. (Use key to be nil if accessing public channel)
