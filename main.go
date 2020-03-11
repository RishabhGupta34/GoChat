package main

import (
	"fmt"
	"net"
	"net/http"

	"GoChat/server"
	//"gecgithub01.walmart.com/r0g03iz/GoChat/server"

	"github.com/gorilla/mux"
)

func main() {
	var ServerVars *server.Variables = &server.Variables{}
	server.Initialise(ServerVars, 2, 1)
	ln, err := net.Listen("tcp", fmt.Sprintf("%s:%s", ServerVars.Conf.Host, ServerVars.Conf.TelPort)) //For establishing a server that listens on the given port
	server.LoggerF(ServerVars, err)
	defer ln.Close()
	r := mux.NewRouter() // Multiplexer for handling the API calls
	r.HandleFunc("/send", ServerVars.Send).Methods("POST")
	r.HandleFunc("/history/{channel}/{key}", ServerVars.History).Methods("GET")
	go func() {
		err := http.ListenAndServe(fmt.Sprintf("%s:%s", ServerVars.Conf.Host, ServerVars.Conf.APIPort), r) // For API calls
		server.LoggerF(ServerVars, err)
	}()
	for {
		conn, err := ln.Accept()
		if err != nil {
			ServerVars.LogError.Println(err) // Logging errors
			continue
		}
		go server.HandleConn(ServerVars, &conn) // A new go-routine for a new client
	}
}
