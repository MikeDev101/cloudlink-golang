package cloudlink

// This structure represents the JSON formatting used for the current CloudLink formatting scheme.
// Values that are not specific to one type are represented with any.
type PacketUPL struct {
	Cmd      string   `json:"cmd"`
	Name     any      `json:"name,omitempty"`
	Val      any      `json:"val"`
	ID       any      `json:"id,omitempty"`
	Rooms    []string `json:"rooms,omitempty"`
	Room     string   `json:"room,omitempty"`
	Listener any      `json:"listener,omitempty"`
	Code     string   `json:"code,omitempty"`
	CodeID   int      `json:"code_id,omitempty"`
}

// This structure represents the JSON formatting the Scratch cloud variable protocol uses.
// Values that are not specific to one type are represented with any.
type Scratch struct {
	Method    string `json:"method"`
	ProjectID string `json:"project_id,omitempty"`
	Username  string `json:"user,omitempty"`
	Value     any    `json:"value"`
	Name      string `json:"name,omitempty"`
	NewName   string `json:"new_name,omitempty"`
}
