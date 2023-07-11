package cloudlink

import (
	// "errors"

	"github.com/bwmarrin/snowflake"
	"github.com/google/uuid"
)

/*
func containsValueClientlist(slice map[snowflake.ID]*Client, value interface{}) bool {
	for _, item := range slice {
		if item == value {
			return true
		}
	}
	return false
}

func containsValue(slice []interface{}, value interface{}) bool {
	for _, item := range slice {
		if item == value {
			return true
		}
	}
	return false
}
*/

// Make a temporary-use copy of a client with RWMutex read locks. This is to safely permit multiple reads from a client.
func (client *Client) TempCopy() *Client {
	// Acquire client read lock
	client.RLock()

	// Copy values
	username := client.username
	protocol := client.protocol
	rooms := TempCopyRooms(client.rooms)
	id := client.id
	uuid := client.uuid
	handshake := client.handshake

	// Free client read lock
	client.RUnlock()

	// Return copied client
	return &Client{
		username:  username,
		handshake: handshake,
		id:        id,
		rooms:     rooms,
		protocol:  protocol,
		uuid:      uuid,
	}
}

// Gathers a map of all Snowflake IDs representing Clients in a Room or Manager.
func GatherSnowflakeIDs(clientstore interface{}) map[snowflake.ID]*Client {
	allids := make(map[snowflake.ID]*Client)
	var readmode uint8
	var tmproom *Room
	var tmpmgr *Manager

	// Type assertions
	switch clientstore.(type) {
	case *Room:
		readmode = 1
	case *Manager:
		readmode = 2
	default:
		return nil
	}

	// Get type, lock, and then read clients
	var clients map[snowflake.ID]*Client
	switch readmode {
	case 1:
		tmproom = clientstore.(*Room)
		tmproom.RLock()
		clients = tmproom.clients
	case 2:
		tmpmgr = clientstore.(*Manager)
		tmpmgr.RLock()
		clients = tmpmgr.clients
	}

	// Gather Snowflake IDs
	for _, client := range clients {

		// Require a set username and a compatible protocol
		tmpclient := client.TempCopy()
		if (tmpclient.username == nil) || (tmpclient.protocol != 1) {
			continue
		}

		allids[client.id] = client
	}

	// Free lock
	switch readmode {
	case 1:
		tmproom.RUnlock()
	case 2:
		tmpmgr.RUnlock()
	}

	// Return collected Snowflake IDs
	return allids
}

// Gathers a map of all UUIDs representing Clients in a Room or Manager.
func GatherUUIDs(clientstore interface{}) map[uuid.UUID]*Client {
	alluuids := make(map[uuid.UUID]*Client)
	var readmode uint8
	var tmproom *Room
	var tmpmgr *Manager

	// Type assertions
	switch clientstore.(type) {
	case *Room:
		readmode = 1
	case *Manager:
		readmode = 2
	default:
		return nil
	}

	// Get type, lock, and then read clients
	var clients map[snowflake.ID]*Client
	switch readmode {
	case 1:
		tmproom = clientstore.(*Room)
		tmproom.RLock()
		clients = tmproom.clients
	case 2:
		tmpmgr = clientstore.(*Manager)
		tmpmgr.RLock()
		clients = tmpmgr.clients
	}

	// Gather UUIds
	for _, client := range clients {

		// Require a set username and a compatible protocol
		tmpclient := client.TempCopy()
		if (tmpclient.username == nil) || (tmpclient.protocol != 1) {
			continue
		}

		alluuids[client.uuid] = client
	}

	// Free lock
	switch readmode {
	case 1:
		tmproom.RUnlock()
	case 2:
		tmpmgr.RUnlock()
	}

	// Return collected UUIDs
	return alluuids
}

// Gathers a map of all UserObjects representing Clients in a Room or Manager.
func GatherUserObjects(clientstore interface{}) map[*UserObject]*Client {
	alluserobjects := make(map[*UserObject]*Client)
	var readmode uint8
	var tmproom *Room
	var tmpmgr *Manager

	// Type assertions
	switch clientstore.(type) {
	case *Room:
		readmode = 1
	case *Manager:
		readmode = 2
	default:
		return nil
	}

	// Get type, lock, and then read clients
	var clients map[snowflake.ID]*Client
	switch readmode {
	case 1:
		tmproom = clientstore.(*Room)
		tmproom.RLock()
		clients = tmproom.clients
	case 2:
		tmpmgr = clientstore.(*Manager)
		tmpmgr.RLock()
		clients = tmpmgr.clients
	}

	// Gather usernames
	for _, client := range clients {

		// Require a set username and a compatible protocol
		tmpclient := client.TempCopy()
		if (tmpclient.username == nil) || (tmpclient.protocol != 1) {
			continue
		}

		alluserobjects[client.GenerateUserObject()] = client
	}

	// Free lock
	switch readmode {
	case 1:
		tmproom.RUnlock()
	case 2:
		tmpmgr.RUnlock()
	}

	// Return collected UserObjects
	return alluserobjects
}

// Gathers a map of all Usernames representing multiple Clients in a Room or Manager.
func GatherUsernames(clientstore interface{}) map[interface{}][]*Client {
	allusernames := make(map[interface{}][]*Client)
	var readmode uint8
	var tmproom *Room
	var tmpmgr *Manager

	// Type assertions
	switch clientstore.(type) {
	case *Room:
		readmode = 1
	case *Manager:
		readmode = 2
	default:
		return nil
	}

	// Get type, lock, and then read clients
	var clients map[snowflake.ID]*Client
	switch readmode {
	case 1:
		tmproom = clientstore.(*Room)
		tmproom.RLock()
		clients = tmproom.clients
	case 2:
		tmpmgr = clientstore.(*Manager)
		tmpmgr.RLock()
		clients = tmpmgr.clients
	}

	// Gather usernames
	for _, client := range clients {

		// Require a set username and a compatible protocol
		tmpclient := client.TempCopy()
		if (tmpclient.username == nil) || (tmpclient.protocol != 1) {
			continue
		}

		allusernames[client.username] = append(allusernames[client.username], client)
	}

	// Free lock
	switch readmode {
	case 1:
		tmproom.RUnlock()
	case 2:
		tmpmgr.RUnlock()
	}

	// Return collected usernames
	return allusernames
}

// Takes a UUID, Snowflake ID, Username, or UserObject query and returns either a single Client (UUID, Snowflake, UserObject) or multiple Clients (username).
func (room *Room) FindClient(query any) interface{} {
	// Interface handling
	switch query.(type) {

	// Snowflake search
	case snowflake.ID:
		snowflakequery := query
		ids := GatherSnowflakeIDs(room)
		if _, ok := ids[snowflakequery.(snowflake.ID)]; ok {
			return ids[snowflakequery.(snowflake.ID)] // Returns *Client
		}

	// UUID search
	case uuid.UUID:
		uuidquery := query
		uuids := GatherUUIDs(room)
		if _, ok := uuids[uuidquery.(uuid.UUID)]; ok {
			return uuids[uuidquery.(uuid.UUID)] // Returns *Client
		}

	// User object search
	case *UserObject:
		userobjectquery := query
		userobjects := GatherUserObjects(room)
		if _, ok := userobjects[userobjectquery.(*UserObject)]; ok {
			return userobjects[userobjectquery.(*UserObject)] // Returns *Client
		}

	// Username search
	default:
		usernames := GatherUsernames(room)
		if _, ok := usernames[query]; ok {
			return usernames[query] // Returns array of *Client
		}
	}

	// Unsupported type
	return nil
}

// Generates a value for client identification.
func (client *Client) GenerateUserObject() *UserObject {
	client.RLock()
	defer client.RUnlock()
	if client.username != nil {
		return &UserObject{
			Id:       client.id,
			Username: client.username,
			Uuid:     client.uuid,
		}
	} else {
		return &UserObject{
			Id:   client.id,
			Uuid: client.uuid,
		}
	}
}

// Gathers all user objects in a room, and generates a userlist.
func (room *Room) GenerateUserList() []*UserObject {
	var objarray []*UserObject

	// Gather all UserObjects
	objstore := GatherUserObjects(room)

	// Convert to array
	for _, client := range objstore {
		objarray = append(objarray, client.GenerateUserObject())
	}

	return objarray
}

// Creates a temporary deep copy of a client's rooms map attribute.
func TempCopyRooms(origin map[interface{}]*Room) map[interface{}]*Room {
	clone := make(map[interface{}]*Room, len(origin))
	for x, y := range origin {
		clone[x] = y
	}
	return clone
}

func RemoveValue(slice []interface{}, indexRemove int) []interface{} {
	// Swap the element to remove with the last element
	slice[indexRemove] = slice[len(slice)-1]

	// Remove the last element
	slice = slice[:len(slice)-1]
	return slice
}

func GetValue(slice []interface{}, target interface{}) int {
	for i, value := range slice {
		if value == target {
			return i
		}
	}
	return -1 // Indicates that the value was not found
}

/*
func appendToSlice(slice []interface{}, elements ...interface{}) ([]interface{}, error) {
	// Use the ellipsis (...) to pass multiple elements as arguments to append
	newSlice := append(slice, elements...)

	// Check if the length of the new slice is as expected
	if len(newSlice) != len(slice)+len(elements) {
		return nil, errors.New("failed to append elements to slice")
	}

	return newSlice, nil
}
*/
