package cloudlink

import (
	"encoding/json"
	"fmt"
	"reflect"
)

// ValidatePacket returns a true or false value based off of if it follows UPL2 or UPL1 depending on the type specified by the user
func ValidatePacket(packet []byte, packetVersion string) bool {
	if packetVersion == "UPL2" {
		var packetstore PacketUPL2
		err := json.Unmarshal(packet, &packetstore)
		if err != nil {
			return false
		}

		if reflect.TypeOf(packetstore.Cmd).Name() != "string" {
			return false
		}
		if reflect.TypeOf(packetstore.Name).Name() != "string" && reflect.TypeOf(packetstore.Name).Name() != "int" {
			return false
		}
		if reflect.TypeOf(packetstore.Code).Name() != "string" {
			return false
		}
		if reflect.TypeOf(packetstore.CodeID).Name() != "int" {
			return false
		}
		return true

	} else if packetVersion == "UPL1" {

		var packetstore PacketUPL1
		err := json.Unmarshal(packet, &packetstore)
		if err != nil {
			return false
		}

		if reflect.TypeOf(packetstore.Cmd).Name() != "string" {
			return false
		}
		if reflect.TypeOf(packetstore.ID).Name() != "string" {
			return false
		}
		return true
	} else {
		fmt.Printf("Error Validating Packet, %s is not a valid packet type", packetVersion)
		return false
	}
}
