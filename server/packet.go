package cloudlink

// This structure represents the JSON formatting used for the legacy CloudLink UPLv1 formatting scheme.
// Values that are not specific to one type are represented with any.
type PacketUPL1 struct {
	Cmd  string `json:"cmd"`
	Val  any    `json:"val"`
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

// This structure represents the JSON formatting used for the current CloudLink UPLv2.1 formatting scheme.
// Values that are not specific to one type are represented with any.
type PacketUPL2 struct {
	Cmd      string `json:"cmd"`
	Name     any    `json:"name,omitempty"`
	Val      any    `json:"val"`
	ID       any    `json:"id,omitempty"`
	Rooms    any    `json:"rooms,omitempty"`
	Listener any    `json:"listener,omitempty"`
	Code     string `json:"code,omitempty"`
	CodeID   int    `json:"code_id,omitempty"`
}

// This structure represents the JSON formatting the Scratch cloud variable protocol uses.
// Values that are not specific to one type are represented with any.
type Scratch struct {
	Method    string `json:"method"`
	ProjectID any    `json:"project_id,omitempty"`
	Username  string `json:"user,omitempty"`
	Value     any    `json:"value"`
	Name      string `json:"name,omitempty"`
	NewName   string `json:"new_name,omitempty"`
}
