package main

type Server struct {
	// Other
	EnableLogs  bool
	IDCounter   int
	IPBlocklist []string

	// Config
	ServerVersion        string
	RejectClients        bool
	EnableScratchSupport bool
	CheckIPAddresses     bool
	EnableMOTD           bool
	MOTDMessage          string

	// Managing methods
	CustomMethods     func()
	DisabledMethods   map[string]struct{}
	MethodCallbacks   map[string]interface{}
	ListenerCallbacks map[string]interface{}
	SafeMethods       map[string]struct{}
}
