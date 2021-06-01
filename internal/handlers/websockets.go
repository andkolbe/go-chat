package handlers

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

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

	// // add a new user to the map when they connect to the server with a websocket connection
	// conn := WebSocketConnection{Conn: ws}
	// clients[conn] = ""

	err = ws.WriteJSON(response)
	if err != nil {
		log.Println(err)
	}

	// // start go routine. Run it forever
	// go ListenForWS(&conn)
}
