package main

var ulist = []interface{}{}

type UlistSub struct {
	Method string
	Val    interface{}
}

func AddUser(name interface{}, listener interface{}) PacketUPL2 {
	if !containsValue(ulist, name) {
		temp, err := appendToSlice(ulist, name)
		if err != nil {
			return PacketUPL2{
				Cmd: "statuscode",
				Val: "E:105 | Internal server error",
			}
		} else {
			ulist = temp
			return GetULIST(listener)
		}
	}
	return GetULIST(listener)
}
func GetULIST(listener interface{}) PacketUPL2 {
	if listener == nil {
		return PacketUPL2{
			Cmd: "ulist",
			Val: UlistSub{
				Method: "set",
				Val:    ulist,
			},
			Rooms: "default",
		}
	} else {
		return PacketUPL2{
			Cmd: "ulist",
			Val: UlistSub{
				Method: "set",
				Val:    ulist,
			},
			Rooms:    "default",
			Listener: listener,
		}
	}
}
