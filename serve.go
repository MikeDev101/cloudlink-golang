package main

import (
	"log"
	"net/http"
)

// You must call initServer beforehand to use this function.

/*
serveWS starts the CloudLink server in insecure (ws://) mode.
You must call initServer beforehand to use this function.
*/
func ServeWS(host string) {
	// Display a startup version string.
	log.Printf("CloudLink Server (Go Edition) v%v - Listening to ws://%v", Version, host)

	// Begin running the websocket server.
	log.Fatal(http.ListenAndServe(host, nil))
}

/*
serveWSS starts the CloudLink server in secure (wss://) mode.
You must have a certificate file (cert, or cert.pem) and
a private key (key, or key.pem) file present and trusted
by your client(s) before running.
*/
func ServeWSS(host string, cert string, key string) {
	// Display a startup version string.
	log.Printf("CloudLink Server (Go Edition) v%v - Listening to wss://%v", Version, host)

	// Begin running the websocket server with SSL support.
	log.Fatal(http.ListenAndServeTLS(host, cert, key, nil))
}
