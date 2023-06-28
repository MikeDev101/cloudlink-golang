package cloudlink

import (
	"fmt"
)

type UserObject struct {
	Id       string `json:"id,omitempty"`
	Username string `json:"username,omitempty"`
	Uuid     string `json:"uuid,omitempty"`
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

// CL4MethodHandler is a method that gets created when a CL-formatted message gets handled by MessageHandler.
func CL4MethodHandler(client *Client, message *PacketUPL) {
	// TODO: finish this
	switch message.Cmd {
	case "handshake":

		client.RLock()
		defer client.RUnlock()

		// Don't re-broadcast this data if the handshake command was already used
		if !client.handshake {
			client.handshake = true

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

			// Send gmsg states
			for _, room := range client.rooms {
				UnicastMessage(client, &PacketUPL{
					Cmd:  "gmsg",
					Val:  room.gmsgState,
					Room: room.name,
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
		// BUG: if the "rooms" value isn't a list/array, the connection fails

		// Argument "rooms" not specified
		if len(message.Rooms) == 0 {

			// Use all subscribed rooms
			for _, room := range client.rooms {

				// Update room gmsg state
				room.gmsgStateMutex.Lock()
				room.gmsgState = message.Val
				room.gmsgStateMutex.Unlock()

				// Broadcast the new state
				room.gmsgStateMutex.RLock()
				MulticastMessage(room.clients, &PacketUPL{
					Cmd:  "gmsg",
					Val:  room.gmsgState,
					Room: room.name,
				})
				room.gmsgStateMutex.RUnlock()
			}

		} else {
			// Use specified rooms
			for _, room := range message.Rooms {

				// Convert input to string
				tmproom := fmt.Sprint(room)

				// Check if room is valid and is subscribed
				if _, ok := client.rooms[tmproom]; ok {
					room := client.rooms[tmproom]

					// Update room gmsg state
					room.gmsgStateMutex.Lock()
					room.gmsgState = message.Val
					room.gmsgStateMutex.Unlock()

					// Broadcast the new state
					room.gmsgStateMutex.RLock()
					MulticastMessage(room.clients, &PacketUPL{
						Cmd:  "gmsg",
						Val:  room.gmsgState,
						Room: room.name,
					})
					room.gmsgStateMutex.RUnlock()
				}
			}
		}

	case "pmsg":
		// Require username to be set before usage
		client.RLock()
		defer client.RUnlock()
		if !client.usernameset {
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
		// Convert input to string
		tmpname := fmt.Sprint(message.Val)

		// Prevent changing username
		client.RLock()
		if client.usernameset {
			UnicastMessage(client, &PacketUPL{
				Cmd:      "statuscode",
				Code:     "E:107 | ID already set",
				CodeID:   107,
				Val:      client.GenerateUserObject(),
				Listener: message.Listener,
			})
			client.RUnlock()
			return
		}
		client.RUnlock()

		// Update client attributes
		client.Lock()
		client.username = tmpname
		client.usernameset = true
		client.Unlock()

		// Send status code
		UnicastMessage(client, &PacketUPL{
			Cmd:      "statuscode",
			Code:     "I:100 | OK",
			CodeID:   100,
			Val:      client.GenerateUserObject(),
			Listener: message.Listener,
		})

	case "gvar":
		// BUG: if the "rooms" value isn't a list/array, the connection fails

		// Argument "rooms" not specified
		if len(message.Rooms) == 0 {

			// Use all subscribed rooms
			for _, room := range client.rooms {

				// Convert input to string
				tmpname := fmt.Sprint(message.Name)

				// Update room gvar state
				room.gvarStateMutex.Lock()
				room.gvarState[tmpname] = message.Val
				room.gvarStateMutex.Unlock()

				// Broadcast the new state
				room.gvarStateMutex.RLock()
				MulticastMessage(room.clients, &PacketUPL{
					Cmd:  "gvar",
					Name: tmpname,
					Val:  room.gvarState[tmpname],
					Room: room.name,
				})
				room.gvarStateMutex.RUnlock()
			}

		} else {
			// Use specified rooms
			for _, room := range message.Rooms {

				// Convert input to string
				tmpname := fmt.Sprint(message.Name)

				// Convert input to string
				tmproom := fmt.Sprint(room)

				// Check if room is valid and is subscribed
				if _, ok := client.rooms[tmproom]; ok {
					room := client.rooms[tmproom]

					// Update room gmsg state
					room.gvarStateMutex.Lock()
					room.gvarState[tmpname] = message.Val
					room.gvarStateMutex.Unlock()

					// Broadcast the new state
					room.gvarStateMutex.RLock()
					MulticastMessage(room.clients, &PacketUPL{
						Cmd:  "gvar",
						Val:  room.gvarState[tmpname],
						Room: room.name,
					})
					room.gvarStateMutex.RUnlock()
				}
			}
		}

	case "pvar":
		// Require username to be set before usage
		client.RLock()
		if !client.usernameset {
			UnicastMessage(client, &PacketUPL{
				Cmd:      "statuscode",
				Code:     "E:111 | ID required",
				CodeID:   111,
				Val:      client.GenerateUserObject(),
				Listener: message.Listener,
			})
			client.RUnlock()
			return
		}
		client.RUnlock()

	case "link":
		// Require username to be set before usage
		client.RLock()
		if !client.usernameset {
			UnicastMessage(client, &PacketUPL{
				Cmd:      "statuscode",
				Code:     "E:111 | ID required",
				CodeID:   111,
				Val:      client.GenerateUserObject(),
				Listener: message.Listener,
			})
			client.RUnlock()
			return
		}
		client.RUnlock()

	case "unlink":
		// Require username to be set before usage
		client.RLock()
		if !client.usernameset {
			UnicastMessage(client, &PacketUPL{
				Cmd:      "statuscode",
				Code:     "E:111 | ID required",
				CodeID:   111,
				Val:      client.GenerateUserObject(),
				Listener: message.Listener,
			})
			client.RUnlock()
			return
		}
		client.RUnlock()

	case "direct":
		// Require username to be set before usage
		client.RLock()
		if !client.usernameset {
			UnicastMessage(client, &PacketUPL{
				Cmd:      "statuscode",
				Code:     "E:111 | ID required",
				CodeID:   111,
				Val:      client.GenerateUserObject(),
				Listener: message.Listener,
			})
			client.RUnlock()
			return
		}
		client.RUnlock()

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