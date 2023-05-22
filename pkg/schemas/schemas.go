package schemas

// CloudlinkRootSchema is the pre-defined formatting for incoming CloudLink messages.
type CloudlinkRootSchema struct {
	Cmd string
	Val any
	Id  any
}
