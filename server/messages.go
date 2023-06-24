package cloudlink

import (
	"log"

	"github.com/goccy/go-json"
	"github.com/gofiber/contrib/websocket"
)

func JSONDump(message any) []byte {
	payload, _ := json.Marshal(message)
	return payload
}

// MulticastMessageRooms takes a Room struct and broadcasts a payload to all clients stored within the room's clients map.
func MulticastMessageRooms(room *Room, message any) {
	for _, client := range room.clients {
		// Spawn goroutines to multicast the payload
		go UnicastMessage(&client, message)
	}
}

// MulticastMessageRooms takes a client manager and broadcasts a payload to all clients stored within the manager's clients map.
func MulticastMessageManager(manager *Manager, message any) {
	for _, client := range manager.clients {
		// Spawn goroutines to multicast the payload
		go UnicastMessage(&client, message)
	}
}

// UnicastMessageAny broadcasts a payload to a singular client.
func UnicastMessage(client *Client, message any) {
	// Echo message back to client - Log errors if any
	if err := client.connection.WriteMessage(websocket.TextMessage, JSONDump(message)); err != nil {
		log.Printf("Client %s (%s) TX error: %s", client.id, client.uuid, err)
	}
}
