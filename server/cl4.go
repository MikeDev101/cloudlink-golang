package cloudlink

import (
	"github.com/bwmarrin/snowflake"
	"github.com/google/uuid"
)

type UserObject struct {
	Id       snowflake.ID `json:"id,omitempty"`
	Username string       `json:"username,omitempty"`
	Uuid     uuid.UUID    `json:"uuid,omitempty"`
}

// Generates a value for client identification.
func GenerateUserObject(client *Client) *UserObject {
	if client.usernameset {
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

// CL4MethodHandler is a method that gets created when a CL-formatted message gets handled by MessageHandler.
func CL4MethodHandler(client *Client, message PacketUPL2) {
	// TODO: finish this
	switch message.Cmd {
	case "handshake":

		// Send the client's IP address
		if client.manager.Config.CheckIPAddresses {
			UnicastMessage(client, &PacketUPL2{
				Cmd: "client_ip",
				Val: client.connection.Conn.RemoteAddr().String(),
			})
		}

		// Send the server version info
		UnicastMessage(client, &PacketUPL2{
			Cmd: "server_version",
			Val: ServerVersion,
		})

		// Send MOTD
		if client.manager.Config.EnableMOTD {
			UnicastMessage(client, &PacketUPL2{
				Cmd: "motd",
				Val: client.manager.Config.MOTDMessage,
			})
		}

		// Send Client's object
		UnicastMessage(client, &PacketUPL2{
			Cmd: "client_obj",
			Val: GenerateUserObject(client),
		})

	case "gmsg":
		break
	case "pmsg":
		break
	case "setid":
		break
	case "gvar":
		break
	case "pvar":
		break
	case "link":
		break
	case "unlink":
		break
	case "direct":
		break
	case "echo":
		UnicastMessage(client, message)
	default:
		break
	}
}
