package main

// Import dependencies.
import (
	"log"
	"net/http"

	// "reflect"
	"sync"

	"github.com/bwmarrin/snowflake"
	"github.com/goccy/go-json"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// Version - Used to identify specific versions of the CloudLink Golang server.
var Version string = "0.1.0"

// Configuration
var MessageOfTheDay = "Cloudlink Golang Edition! This works somehow!"
var CheckIPAddresses bool = true

// WebsocketUpgrader is used to upgrade HTTP(s) requests into websocket connections.
var WebsocketUpgrader = websocket.Upgrader{
	ReadBufferSize:  2048,
	WriteBufferSize: 2048,

	// Check for originating domain. Does nothing at the moment.
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Client is a websocket client.
type Client struct {
	connection *websocket.Conn
	manager    *Manager
	id         snowflake.ID
	uuid       uuid.UUID
}

// ClientList is a map used to keep track of clients.
type ClientList map[*Client]Client

// Manager is used to store all values for each client.
type Manager struct {
	clients ClientList

	// Used to lock states before editing a client.
	sync.RWMutex
}

// NewClient is used to initialize a new Client with all attributes initialized.
func NewClient(conn *websocket.Conn, manager *Manager) *Client {

	// Generate Snowflake ID
	node, err := snowflake.NewNode(1)
	if err != nil {
		log.Fatalln(err, 3)
	}

	return &Client{
		connection: conn,
		manager:    manager,
		id:         node.Generate(),
		uuid:       uuid.New(),
	}
}

// NewManager is used to create a manager struct and initialize its values.
func NewManager() *Manager {
	return &Manager{
		clients: make(ClientList),
	}
}

// AddClient adds a client to the ClientList
func (m *Manager) AddClient(client *Client) {
	// Create a lock to modify values
	m.Lock()
	defer m.Unlock()

	// Add the client
	m.clients[client] = *client

	// Log created client
	log.Printf("Client %v Connected.", client.id)
	log.Printf("There are %v clients connected.", len(m.clients))
}

// RemoveClient removes a client from the ClientList and cleans things up
func (m *Manager) RemoveClient(client *Client) {
	// Create a lock to modify values
	m.Lock()
	defer m.Unlock()

	// Verify if a client exists, and delete if true
	if _, ok := m.clients[client]; ok {
		client.connection.Close()

		// Remove client from ClientList
		delete(m.clients, client)

		// Log deleted client
		log.Printf("Client %v Disconnected.", client.id)
		log.Printf("There are %v clients connected.", len(m.clients))
	}
}

// Server is used to handle new connections, create managers, and process messages.
func (mgr *Manager) Server(w http.ResponseWriter, r *http.Request) {
	// Attempt to upgrade the Websocket connection. If failed, abort and log error.
	con, err := WebsocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	// Create a new client
	client := NewClient(con, mgr)

	// Add client to manager
	mgr.AddClient(client)

	// Begin listening for messages
	client.MessageHandler(mgr)
}

// UnicastPacket will send a message to a single client.
func UnicastPacket(client Client, message PacketUPL2) {

	// Marshal the JSON message
	data, err := json.Marshal(message)
	if err != nil {
		log.Println(err)
		return
	}

	// Transmit
	if err := client.connection.WriteMessage(websocket.TextMessage, data); err != nil {
		log.Println(err)
		return
	}
}

// MulticastMessage will send a message to multiple clients.
func MulticastMessage(clients ClientList, message PacketUPL2) {
	for _, c := range clients {
		// log.Printf("%v (%v)", c.id, c.uuid)
		UnicastPacket(c, message)
	}
}

// Main function for handling websocket messages.
func (c *Client) MessageHandler(mgr *Manager) {

	// Gracefully close the connection once function is complete
	defer c.manager.RemoveClient(c)

	// Connection loop
	for {
		// Receive incoming message
		_, messagePayload, err := c.connection.ReadMessage()

		// Handle errors
		if err != nil {
			break
		}

		// Parse JSON
		var packet PacketUPL2
		parse_err := json.Unmarshal([]byte(messagePayload), &packet)

		// Handle JSON Parsing errors
		if parse_err != nil {
			log.Println("JSON Parsing error:", parse_err)
			continue
		}

		// Display input message after parsing
		// log.Println("InputType:", messageType, ", InputMessage:", InputMessage)

		// TODO: Message command processing code

		switch packet.Cmd {
		case "handshake":
			HandleHandshake(mgr)
		case "gmsg":
			HandleGMSG(mgr, packet.Val)
		default:
			ServerResponse := PacketUPL2{
				Cmd: "direct",
				Val: "I:100 | OK",
			}
			MulticastMessage(mgr.clients, ServerResponse)
		}

	}
}

// initServer initializes the CloudLink server.
func initServer() {
	// Create a new Manager instance to manage websocket instances
	manager := NewManager()

	// Serve the websocket connection at root route: /
	http.HandleFunc("/", manager.Server)
}

// main is the start of the server program.
func main() {
	// Configure the websocket server.
	initServer()

	// Run server
	ServeWS("localhost:3000")
}
