package cloudlink

import (
	"fmt"
)

type UserObject struct {
	Id       string      `json:"id,omitempty"`
	Username interface{} `json:"username,omitempty"`
	Uuid     string      `json:"uuid,omitempty"`
}

func CL4ProtocolDetect(client *Client) {
	if client.protocol == 0 {
		// Update client attributes
		client.Lock()
		client.protocol = 1 // CL4
		client.Unlock()

		// Creates room if it does not exist already
		defaultroom := client.manager.CreateRoom("default")

		// Add the client to the room
		defaultroom.SubscribeClient(client)
	}
}

func (room *Room) BroadcastUserlistEvent(event string, value interface{}, exclude *Client) {
	// Create a dummy manager for selecting clients
	dummy := DummyManager(room.name)

	// Separate compatible clients
	for _, client := range room.clients {
		// Exclude
		if (exclude != nil) && (exclude.id == client.id) {
			continue
		}

		// Require a set username and a compatible protocol
		if (client.username == nil) || (client.protocol != 1) {
			continue
		}

		// Add client if passed
		dummy.AddClient(client)
	}

	// Broadcast state
	MulticastMessage(dummy.clients, &PacketUPL{
		Cmd:   "ulist",
		Val:   value,
		Mode:  event,
		Rooms: room.name,
	})
}

func (room *Room) BroadcastGmsg(value interface{}) {
	// Update room gmsg state
	room.gmsgStateMutex.Lock()
	room.gmsgState = value
	room.gmsgStateMutex.Unlock()

	// Broadcast the new state
	room.gmsgStateMutex.RLock()
	MulticastMessage(room.clients, &PacketUPL{
		Cmd:   "gmsg",
		Val:   room.gmsgState,
		Rooms: room.name,
	})
	room.gmsgStateMutex.RUnlock()
}

func (room *Room) BroadcastGvar(name string, value interface{}) {
	// Update room gmsg state
	room.gvarStateMutex.Lock()
	room.gvarState[name] = value
	room.gvarStateMutex.Unlock()

	// Broadcast the new state
	room.gvarStateMutex.RLock()
	MulticastMessage(room.clients, &PacketUPL{
		Cmd:   "gvar",
		Name:  name,
		Val:   room.gvarState[name],
		Rooms: room.name,
	})
	room.gvarStateMutex.RUnlock()
}

func (client *Client) RequireIDBeingSet(message *PacketUPL) bool {
	client.RLock()
	usernameset := (client.username != nil)
	client.RUnlock()
	if !usernameset {
		UnicastMessage(client, &PacketUPL{
			Cmd:      "statuscode",
			Code:     "E:111 | ID required",
			CodeID:   111,
			Val:      client.GenerateUserObject(),
			Listener: message.Listener,
		})
	}
	return usernameset
}

func (client *Client) HandleIDSet(message *PacketUPL) bool {
	client.RLock()
	usernameset := (client.username != nil)
	client.RUnlock()
	if usernameset {
		UnicastMessage(client, &PacketUPL{
			Cmd:      "statuscode",
			Code:     "E:107 | ID already set",
			CodeID:   107,
			Val:      client.GenerateUserObject(),
			Listener: message.Listener,
		})
	}
	return usernameset
}

// CL4MethodHandler is a method that gets created when a CL-formatted message gets handled by MessageHandler.
func CL4MethodHandler(client *Client, message *PacketUPL) {
	// TODO: finish this
	switch message.Cmd {
	case "handshake":

		// Read attribute
		client.RLock()
		handshakeDone := client.handshake
		client.RUnlock()

		// Don't re-broadcast this data if the handshake command was already used
		if !handshakeDone {

			// Update attribute
			client.Lock()
			client.handshake = true
			client.Unlock()

			// Send the client's IP address
			if client.manager.Config.CheckIPAddresses {
				UnicastMessage(client, &PacketUPL{
					Cmd: "client_ip",
					Val: client.connection.Conn.RemoteAddr().String(),
				})
			}

			// Send the server version info
			UnicastMessage(client, &PacketUPL{
				Cmd: "server_version",
				Val: ServerVersion,
			})

			// Send MOTD
			if client.manager.Config.EnableMOTD {
				UnicastMessage(client, &PacketUPL{
					Cmd: "motd",
					Val: client.manager.Config.MOTDMessage,
				})
			}

			// Send Client's object
			UnicastMessage(client, &PacketUPL{
				Cmd: "client_obj",
				Val: client.GenerateUserObject(),
			})

			// Send gmsg and ulist states
			for _, room := range client.rooms {
				UnicastMessage(client, &PacketUPL{
					Cmd:   "gmsg",
					Val:   room.gmsgState,
					Rooms: room.name,
				})
				UnicastMessage(client, &PacketUPL{
					Cmd:   "ulist",
					Mode:  "set",
					Val:   room.GenerateUserList(),
					Rooms: room.name,
				})
			}
		}

		// Send status code
		UnicastMessage(client, &PacketUPL{
			Cmd:    "statuscode",
			Code:   "I:100 | OK",
			CodeID: 100,
		})

	case "gmsg":
		// Check if required Val argument is provided
		switch message.Val.(type) {
		case nil:
			UnicastMessage(client, &PacketUPL{
				Cmd:     "statuscode",
				Code:    "E:101 | Syntax",
				CodeID:  101,
				Details: "Message missing required val key",
			})
			return
		}

		// Handle multiple types for room
		switch message.Rooms.(type) {

		// Value not specified in message
		case nil:
			// Use all subscribed rooms
			for _, room := range client.rooms {
				room.BroadcastGmsg(message.Val)
			}

		// Multiple rooms
		case []any:
			for _, room := range message.Rooms.([]any) {

				// Convert input to string
				tmproom := fmt.Sprint(room)

				// Check if room is valid and is subscribed
				if _, ok := client.rooms[tmproom]; ok {
					client.rooms[tmproom].BroadcastGmsg(message.Val)
				}
			}

		// Single room
		case any:
			// Convert input to string
			tmproom := fmt.Sprint(message.Rooms)

			// Check if room is valid and is subscribed
			if _, ok := client.rooms[tmproom]; ok {
				client.rooms[tmproom].BroadcastGmsg(message.Val)
			}
		}

	case "pmsg":
		// Require username to be set before usage
		client.RLock()
		usernameset := (client.username != nil)
		client.RUnlock()

		if !usernameset {
			UnicastMessage(client, &PacketUPL{
				Cmd:      "statuscode",
				Code:     "E:111 | ID required",
				CodeID:   111,
				Val:      client.GenerateUserObject(),
				Listener: message.Listener,
			})
			return
		}

	case "setid":
		// Val datatype validation
		switch message.Val.(type) {
		case string:
		case int64:
		case float64:
		case bool:
		default:
			// Send status code
			UnicastMessage(client, &PacketUPL{
				Cmd:      "statuscode",
				Code:     "E:102 | Datatype",
				CodeID:   102,
				Details:  "Username value (val) must be a string, boolean, float, or int",
				Listener: message.Listener,
			})
			return
		}

		// Prevent changing usernames
		if client.HandleIDSet(message) {
			return
		}

		// Update client attributes
		client.Lock()
		client.username = message.Val
		client.Unlock()

		// Use default room
		for _, room := range client.rooms {
			room.BroadcastUserlistEvent("add", client.GenerateUserObject(), client)
			UnicastMessage(client, &PacketUPL{
				Cmd:   "ulist",
				Mode:  "set",
				Val:   room.GenerateUserList(),
				Rooms: room.name,
			})
		}

		// Send status code
		UnicastMessage(client, &PacketUPL{
			Cmd:      "statuscode",
			Code:     "I:100 | OK",
			CodeID:   100,
			Val:      client.GenerateUserObject(),
			Listener: message.Listener,
		})

	case "gvar":
		// Handle multiple types for room
		switch message.Rooms.(type) {

		// Value not specified in message
		case nil:
			// Use all subscribed rooms
			for _, room := range client.rooms {
				room.BroadcastGvar(fmt.Sprint(message.Name), message.Val)
			}

		// Multiple rooms
		case []any:
			// Use specified rooms
			for _, room := range message.Rooms.([]any) {
				// Check if room is valid and is subscribed
				if _, ok := client.rooms[fmt.Sprint(room)]; ok {
					client.rooms[fmt.Sprint(room)].BroadcastGvar(fmt.Sprint(message.Name), message.Val)
				}
			}

		// Single room
		case any:
			// Check if room is valid and is subscribed
			if _, ok := client.rooms[fmt.Sprint(message.Rooms)]; ok {
				client.rooms[fmt.Sprint(message.Rooms)].BroadcastGvar(fmt.Sprint(message.Name), message.Val)
			}
		}

	case "pvar":
		// Require username to be set before usage
		if !client.RequireIDBeingSet(message) {
			return
		}

	case "link":
		// Require username to be set before usage
		if !client.RequireIDBeingSet(message) {
			return
		}

	case "unlink":
		// Require username to be set before usage
		if !client.RequireIDBeingSet(message) {
			return
		}

	case "direct":
		// Require username to be set before usage
		if !client.RequireIDBeingSet(message) {
			return
		}

	case "echo":
		UnicastMessage(client, message)

	default:
		// Handle unknown commands
		UnicastMessage(client, &PacketUPL{
			Cmd:      "statuscode",
			Code:     "E:109 | Invalid command",
			CodeID:   109,
			Val:      client.GenerateUserObject(),
			Listener: message.Listener,
		})
	}
}
