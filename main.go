package main

import (
	"GoChat/server"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	cp := flag.String("configfile", "", "Full path for the config file")
	flag.Parse()
	if len(*cp) == 0 {
		log.Fatal("Please enter the full path for the server's config file")
	}
	var ServerVars *server.Variables = &server.Variables{}
	server.Initialise(ServerVars, *cp, "main") // main is for differentiating b/w main and test

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
	fmt.Println("Telnet server running on ", ServerVars.Conf.Host, ServerVars.Conf.TelPort)
	fmt.Println("API calls listening on ", ServerVars.Conf.Host, ServerVars.Conf.APIPort)

	for {
		conn, err := ln.Accept()
		if err != nil {
			ServerVars.LogError.Println(err) // Logging errors
			continue
		}
		go server.HandleConn(ServerVars, &conn) // A new go-routine for a new client
	}

}
