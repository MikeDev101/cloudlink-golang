package cloudlink

import (
	"log"
	"sync"

	"github.com/bwmarrin/snowflake"
	"github.com/gofiber/contrib/websocket"
	"github.com/google/uuid"
)

var ServerVersion string = "0.1.0-golang"

type Room struct {
	// Subscribed clients to the room
	clients      map[snowflake.ID]*Client
	clientsMutex sync.RWMutex

	// Global message (GMSG) state
	gmsgState      interface{}
	gmsgStateMutex sync.RWMutex

	// Globar variables (GVAR) states
	gvarState      map[string]any
	gvarStateMutex sync.RWMutex

	// Friendly name for room
	name string

	// Locks states before subscribing/unsubscribing clients
	sync.RWMutex
}

type Manager struct {
	// Friendly name for manager
	name string

	// Registered client sessions
	clients      map[snowflake.ID]*Client
	clientsMutex sync.RWMutex

	// Rooms storage
	rooms      map[string]*Room
	roomsMutex sync.RWMutex

	// Configuration settings
	Config struct {
		RejectClients    bool
		CheckIPAddresses bool
		EnableMOTD       bool
		MOTDMessage      string
	}

	// Used for generating Snowflake IDs
	SnowflakeIDNode *snowflake.Node

	// Locks states before registering sessions
	sync.RWMutex
}

// NewClient assigns a UUID and Snowflake ID to a websocket client, and returns a initialized Client struct for use with a manager's AddClient.
func NewClient(conn *websocket.Conn, manager *Manager) *Client {
	// Request and create a lock before generating ID values
	manager.clientsMutex.Lock()

	// Generate client ID values
	client_id := manager.SnowflakeIDNode.Generate()
	client_uuid := uuid.New()

	// Release the lock
	manager.clientsMutex.Unlock()

	return &Client{
		connection: conn,
		manager:    manager,
		id:         client_id,
		uuid:       client_uuid,
		rooms:      make(map[string]*Room),
		handshake:  false,
	}
}

// Dummy Managers function identically to a normal manager. However, they are used for selecting specific clients to multicast to.
func DummyManager() *Manager {
	return &Manager{
		clients: make(map[snowflake.ID]*Client),
		rooms:   make(map[string]*Room),
	}
}

func New(name string) *Manager {
	node, err := snowflake.NewNode(1)
	if err != nil {
		log.Fatalln(err, 3)
	}

	manager := &Manager{
		name:            name,
		clients:         make(map[snowflake.ID]*Client),
		rooms:           make(map[string]*Room),
		SnowflakeIDNode: node,
	}

	return manager
}

func (manager *Manager) CreateRoom(name string) *Room {
	manager.roomsMutex.RLock()

	// Access rooms map
	_, exists := manager.rooms[name]

	manager.roomsMutex.RUnlock()

	if !exists {
		manager.roomsMutex.Lock()

		log.Printf("[%s] Creating room %s", manager.name, name)

		// Create and prepare the room state
		manager.rooms[name] = &Room{
			name:      name,
			clients:   make(map[snowflake.ID]*Client, 1),
			gmsgState: "",
			gvarState: make(map[string]any),
		}

		manager.roomsMutex.Unlock()
	}

	// Return the room even if it already exists
	return manager.rooms[name]
}

func (room *Room) SubscribeClient(client *Client) {
	room.clientsMutex.Lock()

	// Add client
	room.clients[client.id] = client

	room.clientsMutex.Unlock()
	client.Lock()

	// Add pointer to subscribed room in client's state
	client.rooms[room.name] = room

	client.Unlock()
}

func (room *Room) UnsubscribeClient(client *Client) {
	room.clientsMutex.Lock()

	// Remove client
	delete(room.clients, client.id)

	room.clientsMutex.Unlock()
	client.Lock()

	// Remove pointer to subscribed room from client's state
	delete(client.rooms, room.name)

	client.Unlock()
}

func (manager *Manager) DeleteRoom(name string) {
	manager.roomsMutex.Lock()

	log.Printf("[%s] Destroying room %s", manager.name, name)

	// Delete room
	delete(manager.rooms, name)

	manager.roomsMutex.Unlock()
}

func (manager *Manager) AddClient(client *Client) {
	manager.clientsMutex.Lock()

	// Add client
	manager.clients[client.id] = client

	manager.clientsMutex.Unlock()
}

func (manager *Manager) RemoveClient(client *Client) {
	manager.clientsMutex.Lock()

	// Remove client from manager's clients map
	delete(manager.clients, client.id)

	// Unsubscribe from all rooms and free memory by clearing out empty rooms
	for _, room := range TempCopyRooms(client.rooms) {
		room.UnsubscribeClient(client)

		// Destroy room if empty
		if len(room.clients) == 0 {
			manager.DeleteRoom(room.name)
		}
	}
	manager.clientsMutex.Unlock()
}