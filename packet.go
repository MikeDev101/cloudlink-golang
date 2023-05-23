package main

type PacketUPL1 struct {
	Cmd  string      `json:"cmd"`
	Val  interface{} `json:"val"`
	ID   string      `json:"id"`
	Name string      `json:"name"`
}

type PacketUPL2 struct {
	Cmd      string      `json:"cmd"`
	Name     interface{} `json:"name"`
	Val      interface{} `json:"val"`
	ID       interface{} `json:"id"`
	Rooms    interface{} `json:"rooms"`
	Listener interface{} `json:"listener"`
	Code     string      `json:"code"`
	CodeID   int         `json:"code_id"`
}
