package server

// Contains the functions that handle the HTTP REST API calls(GET,POST)

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

// History is Function for handling HTTP REST GET API calls
func (ServerVars *Variables) History(w http.ResponseWriter, r *http.Request) {
	c := strings.ToLower(mux.Vars(r)["channel"]) // Because channel name is not case sensitive
	k := strings.TrimSpace(mux.Vars(r)["key"])
	ReadMsgAPI(ServerVars, c, k, &w, time.Now())
}

// Send is Function for handling HTTP REST POST API calls
func (ServerVars *Variables) Send(w http.ResponseWriter, r *http.Request) {
	inp, err := ioutil.ReadAll(r.Body)
	LoggerE(ServerVars, err)
	var d DataAPI
	err = json.Unmarshal(inp, &d)
	LoggerE(ServerVars, err)
	d.Channel = strings.ToLower(d.Channel) // Because channel name is not case sensitive
	d.Key = strings.TrimSpace(d.Key)
	SendMsgAPI(ServerVars, d, &w, time.Now())
}
