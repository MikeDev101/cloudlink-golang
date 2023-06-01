package main

var gmsg = []interface{}{}

func AddGMSG(msg interface{}) PacketUPL2 {
	if !containsValue(gmsg, msg) {
		temp, err := appendToSlice(gmsg, msg)
		if err != nil {
			return PacketUPL2{
				Cmd: "statuscode",
				Val: "E:105 | Internal server error",
			}
		} else {
			gmsg = temp
			return PacketUPL2{
				Cmd: "gmsg",
				Val: gmsg,
			}
		}
	}
	return PacketUPL2{
		Cmd:   "gmsg",
		Val:   gmsg,
		Rooms: "default",
	}
}
func GetGMSG() PacketUPL2 {
	temp := len(gmsg)
	if temp > 0 {
		return PacketUPL2{
			Cmd: "gmsg",
			Val: gmsg[temp],
		}
	} else {
		return PacketUPL2{
			Cmd: "gmsg",
			Val: gmsg,
		}
	}
}
