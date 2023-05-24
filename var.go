package main

type Packet_Var struct {
	Name  interface{}
	Value interface{}
}

var gvar []Packet_Var

func AddGVAR(name interface{}, value interface{}) PacketUPL2 {
	temp_packet := Packet_Var{
		Name:  name,
		Value: value,
	}
	gvar = append(gvar, temp_packet)
	return PacketUPL2{
		Cmd:  "gvar",
		Name: temp_packet.Name,
		Val:  temp_packet.Value,
	}
}
