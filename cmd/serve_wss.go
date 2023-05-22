package main

// Import CloudLink
import (
	"github.com/mikedev101/cloudlink-golang/pkg/server"
)

// main is the start of the server program.
func main() {
	// Configure the websocket server.
	server.Init()

	// Run server
	server.ServeWSS("0.0.0.0:8000", "cert.pem", "key.pem")
}