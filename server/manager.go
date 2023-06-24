package cloudlink

import (
	"log"
	"sync"

	"github.com/bwmarrin/snowflake"
	"github.com/gofiber/contrib/websocket"
	"github.com/google/uuid"
)

var ServerVersion string = "0.1.0"

type Room struct {
	// Subscribed clients to the room
	clients map[uuid.UUID]Client

	// Friendly name for room
	name string

	// Locks states before subscribing/unsubscribing clients
	sync.RWMutex
}

type Manager struct {
	// Friendly name for manager
	name string

	// Registered client sessions
	clients map[snowflake.ID]Client

	// Used to avoid race conditions when accessing the clients map
	clientsMutex sync.RWMutex

	// Rooms storage
	rooms map[string]*Room

	// Used to avoid race conditions when accessing the rooms map
	roomsMutex sync.RWMutex

	// Configuration settings
	Config struct {
		EnableLogs       bool
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
	manager.Lock()

	// Generate client ID values
	client_id := manager.SnowflakeIDNode.Generate()
	client_uuid := uuid.New()

	// Release the lock
	manager.Unlock()

	return &Client{
		connection: conn,
		manager:    manager,
		id:         client_id,
		uuid:       client_uuid,
	}
}

// Dummy Managers function identically to a normal manager. However, they are used for selecting specific clients to multicast to.
func DummyManager() *Manager {
	return &Manager{
		clients: make(map[snowflake.ID]Client),
		rooms:   make(map[string]*Room),
	}
}

func New(name string) *Manager {
	node, err := snowflake.NewNode(1)
	if err != nil {
		log.Fatalln(err, 3)
	}

	return &Manager{
		name:            name,
		clients:         make(map[snowflake.ID]Client),
		rooms:           make(map[string]*Room),
		SnowflakeIDNode: node,
	}
}

func (manager *Manager) CreateRoom(name string) {
	// Request and create a lock before modifying values
	manager.roomsMutex.Lock()

	// Create room
	manager.rooms[name] = new(Room)

	// Prepare the room state
	manager.rooms[name].name = name
	manager.rooms[name].clients = make(map[uuid.UUID]Client)

	// Release the lock
	manager.roomsMutex.Unlock()
}

func (manager *Manager) DeleteRoom(name string) {
	// Acquire read lock on the rooms map
	manager.roomsMutex.RLock()

	// Access rooms map
	_, ok := manager.rooms[name]

	// Free the read lock on the rooms map
	manager.roomsMutex.RUnlock()

	if ok {
		// Request and create a lock before modifying values
		manager.roomsMutex.Lock()

		// Delete room
		delete(manager.rooms, name)

		// Release the lock
		manager.roomsMutex.Unlock()
	}
}

func (manager *Manager) AddClient(c *Client) {
	log.Printf("[%s] Client connected: %s (%s)", manager.name, c.id, c.uuid)

	// Lock access to the clients map
	manager.clientsMutex.Lock()

	// Add client
	manager.clients[c.id] = *c

	// Free the lock on the clients map
	manager.clientsMutex.Unlock()
}

func (manager *Manager) RemoveClient(c *Client) {
	// Acquire read lock on the clients map
	manager.clientsMutex.RLock()

	// Access clients map
	_, ok := manager.clients[c.id]

	// Free the read lock on the clients map
	manager.clientsMutex.RUnlock()

	if ok {
		log.Printf("[%s] Client disconnected: %s (%s)", manager.name, c.id, c.uuid)

		// Close the connection
		c.connection.Close()

		// Lock access to the clients map
		manager.clientsMutex.Lock()

		// Delete client session
		delete(manager.clients, c.id)

		// Free the lock on the clients map
		manager.clientsMutex.Unlock()
	}
}
