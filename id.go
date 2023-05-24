package main

var ulist = []interface{}{}

type UlistSub struct {
	Method string
	Val    interface{}
}

func AddUser(name interface{}) PacketUPL2 {
	if !containsValue(ulist, name) {
		temp, err := appendToSlice(ulist, name)
		if err != nil {
			return PacketUPL2{
				Cmd: "statuscode",
				Val: "E:105 | Internal server error",
			}
		} else {
			ulist = temp
			return GetULIST()
		}
	}
	return GetULIST()
}
func GetULIST() PacketUPL2 {
	return PacketUPL2{
		Cmd: "ulist",
		Val: UlistSub{
			Method: "set",
			Val:    ulist,
		},
		Rooms: "default",
	}
}
