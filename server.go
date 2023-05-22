package main

type Server struct {
	EnableLogs           bool
	IDCounter            int
	IPBlocklist          []string
	RejectClients        bool
	EnableScratchSupport bool
	CheckIPAddresses     bool
	EnableMotd           bool
	MotdMessage          string
}
