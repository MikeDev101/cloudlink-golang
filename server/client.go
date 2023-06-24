package cloudlink

import (
	"github.com/bwmarrin/snowflake"
	"github.com/gofiber/contrib/websocket"
	"github.com/google/uuid"
)

// The client struct serves as a template for handling websocket sessions. It stores a client's UUID, Snowflake ID, manager and websocket connection pointer(s).
type Client struct {
	connection  *websocket.Conn
	manager     *Manager
	id          snowflake.ID
	uuid        uuid.UUID
	username    string
	usernameset bool
}
