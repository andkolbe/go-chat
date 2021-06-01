package handlers

import (
	"fmt"
	"log"
	"net/http"
	"sort"

	"github.com/gorilla/websocket"
)

// channel
var wsChan = make(chan WSPayload) // our channel will only accept data of type WsPayload

// place to hold all of the connected users
var clients = make(map[WebSocketConnection]string)

// a wrapper for our websocket
type WebSocketConnection struct {
	*websocket.Conn
}

// the kind of information we are sending to the server
type WSPayload struct {
	Action   string              `json:"action"` // what we expect the server to do
	Username string              `json:"username"`
	Message  string              `json:"message"`
	Conn     WebSocketConnection `json:"-"` // leave it out of the json
}

// all the fields we will be sending back from websockets
type WSJSONResponse struct {
	Action         string   `json:"action"`
	Message        string   `json:"message"`
	MessageType    string   `json:"message_type"`
	ConnectedUsers []string `json:"connected_users"`
}

// use this variable to upgrade to websocket connection
var upgradeConnection = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// upgrades connection to websockets and sends back a JSON response
func WsEndPoint(w http.ResponseWriter, r *http.Request) {
	ws, err := upgradeConnection.Upgrade(w, r, nil) // we aren't going to worry about the response header right now
	if err != nil {
		log.Println(err)
	}

	log.Println("Client connected to endpoint")

	// send back a response back to the client using JSON so it is easy to parse
	var response WSJSONResponse
	response.Message = `<em><small>Connected to Server</small></em>`

	// add a new user to the map when they connect to the server with a websocket connection
	conn := WebSocketConnection{Conn: ws}
	clients[conn] = ""

	err = ws.WriteJSON(response)
	if err != nil {
		log.Println(err)
	}

	// start go routine. Run it forever
	go ListenForWS(&conn)
}

// take the users away from the WsEndPoint and into a go routine
func ListenForWS(conn *WebSocketConnection) {
	// if the go routine stops for any reason, automatically start it back up
	defer func() {
		if r := recover(); r != nil {
			log.Println("Error", fmt.Sprintf("%v", r))
		}
	}()

	var payload WSPayload

	for { // an infinite for loop
		err := conn.ReadJSON(&payload) // listen for an incoming payload from a user
		if err != nil {
			// do nothing. There is no payload
		} else {
			payload.Conn = *conn
			wsChan <- payload // send the contents of the payload to the channel
		}
	}
}

//
func ListenToWsChannel() { // start the channel in the main function
	var response WSJSONResponse // what we are going to send back to the user

	for { // an infinite for loop
		e := <-wsChan // every time we get a payload from the channel, store it in the variable e

		switch e.Action {
		case "username": // if the Action on the event is "username"
			// get a list of user and send it back via broadcast
			clients[e.Conn] = e.Username
			users := getUserList()
			// send the users list back to the client
			response.Action = "list_users"
			response.ConnectedUsers = users
			broadcastToAll(response)

		// send an action to all users to tell them that a user has left
		case "left":
			response.Action = "list_users"
			delete(clients, e.Conn)
			users := getUserList()
			response.ConnectedUsers = users
			broadcastToAll(response)

		// broadcast a message from one client back to all client
		case "broadcast":
			response.Action = "broadcast"
			response.Message = fmt.Sprintf("<strong>%s</strong>: %s", e.Username, e.Message)
			broadcastToAll(response)
		}

	}
}

func getUserList() []string {
	var userList []string
	for _, x := range clients {
		if x != "" { // only display a user on the user list if they have typed in a username
			userList = append(userList, x) // append x onto the slice of users
		}
	}
	// sort the list of users
	sort.Strings(userList)
	return userList
}

// broadcast to all users
func broadcastToAll(response WSJSONResponse) {
	for client := range clients { // loop through all of the clients
		err := client.WriteJSON(response) // write the JSON response to all of the connected clients (users)
		if err != nil {
			log.Println("websocket err")
			_ = client.Close()      // close that connection for that client
			delete(clients, client) // remove them from the map
		}
	}
}
