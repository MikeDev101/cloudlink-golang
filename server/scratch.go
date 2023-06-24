package cloudlink

// ScratchMethodHandler is a method that gets created when a Scratch-formatted message gets handled by MessageHandler.
func ScratchMethodHandler(client *Client, message Scratch) {
	// TODO: finish this
	switch message.Method {
	case "handshake":
		break
	case "set":
		break
	case "create":
		break
	case "rename":
		break
	case "delete":
		break
	default:
		break
	}
}
