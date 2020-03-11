package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/acarl005/stripansi"
	"github.com/gorilla/mux"
)

var w1 sync.WaitGroup
var w2 sync.WaitGroup
var c1 *net.Conn
var c2 *net.Conn

type data struct {
	inp, exp, msg string
}

// WriteMatch is For sending the inp command through client and matching the exp(expected output) with the client read
func WriteMatch(t *testing.T, conn *net.Conn, inp string, exp string, msg string) {
	_, err := (*conn).Write([]byte(inp)) // Sending the inp command
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(100 * time.Millisecond)    // Sleep because the output from the server takes some time
	out := stripansi.Strip(Read(t, conn)) // Reading the client
	if out != exp {
		t.Fatalf("%s Expected:%s Got:%s", msg, exp, out) // Printing fatal message
	}
}

// WriteMatchAPI is For sending the data through POST method and matching it to exp(expected output)
func WriteMatchAPI(t *testing.T, inpC string, inpK string, inpM string, exp string, msg string) {
	inp, err := json.Marshal(map[string]string{"channel": inpC, "key": inpK, "msg": inpM}) // Converting the json data to bytes
	if err != nil {
		t.Fatal(err)
	}
	resp, err := http.Post("http://localhost:8001/send", "application/json", bytes.NewBuffer(inp)) // POST request
	if err != nil {
		t.Fatal(err)
	}
	out, err := ioutil.ReadAll(resp.Body) // Reading the response of the POST request
	if err != nil {
		t.Fatal(err)
	}
	if string(out) != exp {
		t.Fatalf("%s Expected:%s Got:%s", msg, exp, out)
	}
}

// ReadMatchAPI is For performing a GET method and matching it to exp(expected output)
func ReadMatchAPI(t *testing.T, inp string, exp string, msg string) {
	resp, err := http.Get(inp) // GET request
	if err != nil {
		t.Fatal(err)
	}
	out, err := ioutil.ReadAll(resp.Body) // Reading the response of the GET request
	if err != nil {
		t.Fatal(err)
	}
	if string(out) != exp {
		t.Fatalf("%s Expected:%s Got:%s", msg, exp, out)
	}
}

// Read is For reading the data from client's buffer
func Read(t *testing.T, conn *net.Conn) string {
	buff := make([]byte, 2000)
	n, err := (*conn).Read(buff) // Reading input from the user
	if err != nil {
		t.Fatal(err)
	}
	comm := string(buff[:n])
	comm = strings.Trim(comm, "\r\n") // Trimming new line characters
	comm = strings.TrimSpace(comm)    // Trimming extra space
	return comm
}

func TestInit(t *testing.T) {
	var ServerVars *Variables = &Variables{}
	Initialise(ServerVars, 4, 3)
	w1.Add(1) // So that the func Test_init and the below go-routine completes before other test functions are executed
	w2.Add(1) // For waiting the go-routine till the telnet server is online
	go func(t *testing.T) {
		w2.Wait() // For waiting the go-routine till the telnet server is online
		defer w1.Done()
		conn, err := net.Dial("tcp", "localhost:8000")
		if err != nil {
			t.Fatal(err)
		}
		c1 = &conn // Client1 with username: testuser
		// go ReadA(t, con)
		_ = Read(t, c1) // Reading client buffer
		_, err = (*c1).Write([]byte("testuser\n"))
		if err != nil {
			t.Fatal(err)
		}
		time.Sleep(100 * time.Millisecond) // Sleep because the output from the server takes some time
		_ = Read(t, c1)                    // Reading client buffer
		c, err := net.Dial("tcp", "localhost:8000")
		if err != nil {
			t.Fatal(err)
		}
		c2 = &c                                // Client2 with username: testuser2
		_ = Read(t, c2)                        // Reading client buffer
		_, err = c.Write([]byte("testuser\n")) // For testing if same username is accepted or not
		if err != nil {
			t.Fatal(err)
		}
		time.Sleep(100 * time.Millisecond) // Sleep because the output from the server takes some time
		_ = Read(t, c2)                    // Reading client buffer
		_, err = c.Write([]byte("testuser2\n"))
		if err != nil {
			t.Fatal(err)
		}
		time.Sleep(100 * time.Millisecond) // Sleep because the output from the server takes some time
		_ = Read(t, c2)                    // Reading client buffer
	}(t)

	ln, err := net.Listen("tcp", "localhost:8000") //For establishing a server that listens on the port 8000
	if err != nil {
		t.Fatal(err)
	}
	r := mux.NewRouter() // Multiplexer for handling the API calls
	r.HandleFunc("/send", ServerVars.Send).Methods("POST")
	r.HandleFunc("/history/{channel}/{key}", ServerVars.History).Methods("GET")
	go func() {
		err := http.ListenAndServe(fmt.Sprintf("localhost:8001"), r) // For API calls
		if err != nil {
			t.Fatal(err)
		}
	}()
	defer ln.Close()
	w2.Done()
	for i := 0; i < 2; i++ { // Because only 2 clients are connected to the server for testing
		c, err := ln.Accept()
		if err != nil {
			t.Fatal(err) // Logging errors
			i--
			continue
		}
		go HandleConn(ServerVars, &c) // A new go-routine for a new client
	}
}

// Testing HandleConn
func TestHandleConn(t *testing.T) {
	w1.Wait()
	tc := []data{
		{"\\abcd\n", "Wrong option!\n", "Server accepted the wrong command"},
		{"\\channels\n", "all", "Server didn't print channels"},
		{"\\blocklist\n", "You have not blocked anyone\n", "Server didn't print blocklist"},
	}
	for _, tt := range tc {
		WriteMatch(t, c1, tt.inp, tt.exp, tt.msg)
	}
}

// Testing NewChannel
func TestNewChannel(t *testing.T) {
	tc := []data{
		{"\\new gen public\n", "Wrong option!\n", "Server accepted the wrong command"},
		{"\\new gen public nil\n", "New  public  channel gen created", "Not able to create public channel"},
		{"\\new tech private ab12\n", "New  private  channel tech created", "Not able to create private channel"},
		{"\\new all public nil\n", "Channel exists with the name: all", "Server created a new channel with the same name"},
		{"\\new gen1 publ nil\n", "Invalid access type\n", "Server accepted wrong access type"},
	}
	for _, tt := range tc {
		WriteMatch(t, c1, tt.inp, tt.exp, tt.msg)
	}
}

// Testing LeaveChannel
func TestLeaveChannel(t *testing.T) {
	tc := []data{
		{"\\leave\n", "Wrong option!\n", "Server accepted the wrong command"},
		{"\\leave tech\n", "You have left tech", "Not able to leave public channel"},
		{"\\leave gen\n", "You have left gen", "Not able to leave private channel"},
		{"\\leave all\n", "You cannot leave the channel all", "Left the channel all"},
		{"\\leave tech\n", "You are not a member of the channel tech", "Leaving a unsubscribed channel"},
		{"\\leave specs\n", "You are not a member of the channel specs", "Leaving a channel that doesn't exist"},
	}
	for _, tt := range tc {
		WriteMatch(t, c1, tt.inp, tt.exp, tt.msg)
	}
}

// Testing JoinChannel
func TestJoinChannel(t *testing.T) {
	tc := []data{
		{"\\join all\n", "Wrong option!\n", "Server accepted the wrong command"},
		{"\\join all nil\n", "You have already joined the channel all", "Joining already joined channel"},
		{"\\join specs nil\n", "No channel exists with the name: specs", "Joining channel that doesn't exist"},
		{"\\join tech nil\n", "Please enter valid key for joining a private channel\n", "Joining a channel with invalid key"},
		{"\\join tech ab12\n", "You were added to channel tech", "Joining private channel"},
	}
	for _, tt := range tc {
		WriteMatch(t, c1, tt.inp, tt.exp, tt.msg)
	}
}

// Testing Info
func TestInfo(t *testing.T) {
	tc := []data{
		{"\\info\n", "Wrong option!\n", "Server accepted the wrong command"},
		{"\\info specs\n", "You have not joined the channel specs", "Printing info of a channel that doesn't exist"},
		{"\\info gen\n", "You have not joined the channel gen", "Printing info of an unsubscribed channel"},
	}
	for _, tt := range tc {
		WriteMatch(t, c1, tt.inp, tt.exp, tt.msg)
	}
}

// Testing SendMsgC
func TestSendMsgC(t *testing.T) {
	tc := []data{
		{"\\sendc gen\n", "Wrong option!\n", "Server accepted the wrong command"},
		{"\\sendc specs hello\n", "You have not joined the channel specs", "Sending message to a channel that doesn't exist"},
		{"\\sendc gen hello\n", "You have not joined the channel gen", "Sending message to an unsubscribed channel"},
		{"\\sendc all hello\n", "Message \"hello\" sent", "Sending message to a channel"},
	}
	for _, tt := range tc {
		WriteMatch(t, c1, tt.inp, tt.exp, tt.msg)
	}
	_ = Read(t, c2)
	_, err := (*c2).Write([]byte("\\block testuser"))
	if err != nil {
		t.Fatal(err)
	}
	_ = Read(t, c2)
	WriteMatch(t, c1, "\\sendc all hello\n", "Message \"hello\" sent", "Send message to a channel")
}

// Testing BlockUser
func TestBlockUser(t *testing.T) {
	tc := []data{
		{"\\block\n", "Wrong option!\n", "Wrong option"},
		{"\\block testuser1\n", "User doesn't exist", "User doesn't exist"},
		{"\\block testuser\n", "You cannot block yourself", "Blocking yourself"},
		{"\\block testuser2\n", "Blocked user testuser2", "Block user"},
		{"\\block testuser2\n", "User testuser2 already belongs in your blocklist", "Blocking blocked user"},
		{"\\blocklist\n", "testuser2", "Printing blocklist"},
	}
	for _, tt := range tc {
		WriteMatch(t, c1, tt.inp, tt.exp, tt.msg)
	}
}

// Testing UnblockUser
func TestUnblockUser(t *testing.T) {
	tc := []data{
		{"\\unblock\n", "Wrong option!\n", "Wrong option"},
		{"\\unblock testuser1\n", "User doesn't exist", "User doesn't exist"},
		{"\\unblock testuser2\n", "Unblocked user testuser2", "Unblock User"},
		{"\\unblock testuser2\n", "User testuser2 doesn't belong in your blocklist", "Unblock unblocked User"},
	}
	for _, tt := range tc {
		WriteMatch(t, c1, tt.inp, tt.exp, tt.msg)
	}
}

// Testing SendMsgU
func TestSendMsgU(t *testing.T) {
	tc := []data{
		{"\\sendu\n", "Wrong option!\n", "Wrong option"},
		{"\\sendu testuser1 hello\n", "User doesn't exist\n", "User doesn't exist"},
		{"\\sendu testuser2 hello\n", "Message \"hello\" sent", "Send Message User"},
	}
	for _, tt := range tc {
		WriteMatch(t, c1, tt.inp, tt.exp, tt.msg)
	}
	_, err := (*c2).Write([]byte("\\unblock testuser"))
	if err != nil {
		t.Fatal(err)
	}
	_ = Read(t, c2)
	WriteMatch(t, c1, "\\sendu testuser2 hello\n", "Message \"hello\" sent", "Send Message User")
	_ = Read(t, c2)
}

// Testing SendMsgAPI
func TestSendMsgAPI(t *testing.T) {
	WriteMatchAPI(t, "tech", "ab12", "Hello Hi", "Message sent\n", "Private channel sending message")
	_ = Read(t, c1)
	WriteMatchAPI(t, "specs", "nil", "Hello Hi", "No channel exists with the name specs\n", "Channel doesn't exist")
	WriteMatchAPI(t, "tech", "nil", "Hello Hi", "Key is not valid for the private channel tech\n", "Key not valid")
}

// Testing ReadMsgAPI
func TestReadMsgAPI(t *testing.T) {
	ReadMatchAPI(t, "http://localhost:8001/history/tech/nil", "Key is not valid for the private channel tech\n", "Key not valid")
	ReadMatchAPI(t, "http://localhost:8001/history/specs/nil", "No channel exists with the name specs\n", "No channel exists")
}

// Testing DeleteUser
func TestDeleteUser(t *testing.T) {
	_, err := (*c2).Write([]byte("\\block testuser")) // For testing if testuser is removed from the blocklist of testuser2 after testuser has closed the conn.
	if err != nil {
		t.Fatal(err)
	}
	WriteMatch(t, c1, "\\close", "Thanks for using GoChat testuser", "Closing connection")
}
