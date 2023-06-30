package cloudlink

import (
	"fmt"

	"github.com/gofiber/contrib/websocket"
)

func ScratchProtocolDetect(client *Client) {
	if client.protocol == 0 {
		// Update client attributes
		client.Lock()
		client.protocol = 2 // Scratch
		client.Unlock()
	}
}

// ScratchMethodHandler is a method that gets created when a Scratch-formatted message gets handled by MessageHandler.
func ScratchMethodHandler(client *Client, message *Scratch) {
	switch message.Method {
	case "handshake":

		// Validate datatype of project ID
		switch message.ProjectID.(type) {
		case string:
		case bool:
		case int64:
		case float64:
		default:
			client.CloseWithMessage(websocket.CloseUnsupportedData, "Invalid Project ID datatype")
			return
		}

		// Update client attributes
		client.username = message.Username

		// Creates room if it does not exist already
		room := client.manager.CreateRoom(message.ProjectID)

		// Add the client to the room
		room.SubscribeClient(client)

	case "set":
		for _, room := range client.rooms { // Should only ever have 1 entry

			// Convert input to string
			tmpname := fmt.Sprint(message.Name)

			// Update room gvar state
			room.gvarStateMutex.Lock()
			room.gvarState[tmpname] = message.Value
			room.gvarStateMutex.Unlock()

			// Broadcast the new state
			room.gvarStateMutex.RLock()
			MulticastMessage(room.clients, &Scratch{
				Method: "set",
				Value:  room.gvarState[tmpname],
				Name:   tmpname,
			})
			room.gvarStateMutex.RUnlock()
		}

	case "create":
		for _, room := range client.rooms { // Should only ever have 1 entry

			// Convert input to string
			tmpname := fmt.Sprint(message.Name)

			// Update room gvar state
			room.gvarStateMutex.Lock()
			room.gvarState[tmpname] = message.Value
			room.gvarStateMutex.Unlock()

			// Broadcast the new state
			room.gvarStateMutex.RLock()
			MulticastMessage(room.clients, &Scratch{
				Method: "create",
				Value:  room.gvarState[tmpname],
				Name:   tmpname,
			})
			room.gvarStateMutex.RUnlock()
		}

	case "rename":
		for _, room := range client.rooms { // Should only ever have 1 entry

			// Convert inputs to string
			tmpname := fmt.Sprint(message.Name)
			tmpnewname := fmt.Sprint(message.NewName)

			// Retrive old value
			room.gvarStateMutex.RLock()
			oldvalue := room.gvarState[tmpname]
			room.gvarStateMutex.RUnlock()

			// Destroy old value and make a new value
			room.gvarStateMutex.Lock()
			delete(room.gvarState, tmpname)
			room.gvarState[tmpnewname] = oldvalue
			room.gvarStateMutex.Unlock()

			// Broadcast the new state
			MulticastMessage(room.clients, &Scratch{
				Method:  "rename",
				NewName: tmpnewname,
				Name:    tmpname,
			})
		}

	case "delete":
		for _, room := range client.rooms { // Should only ever have 1 entry

			// Convert input to string
			tmpname := fmt.Sprint(message.Name)

			// Destroy value
			room.gvarStateMutex.Lock()
			delete(room.gvarState, tmpname)
			room.gvarStateMutex.Unlock()

			// Broadcast the new state
			MulticastMessage(room.clients, &Scratch{
				Method: "delete",
				Name:   tmpname,
			})
		}

	default:
		break
	}
}
