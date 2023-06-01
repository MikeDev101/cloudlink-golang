package main

import "log"

func HandleHandshake(manager *Manager) {

	ServerResponse := PacketUPL2{
		Cmd: "client_ip",
		Val: "0.0.0.0",
	}
	MulticastMessage(manager.clients, ServerResponse)

	ServerResponse = PacketUPL2{
		Cmd: "server_version",
		Val: Version,
	}
	MulticastMessage(manager.clients, ServerResponse)

	ServerResponse = PacketUPL2{
		Cmd: "motd",
		Val: MessageOfTheDay,
	}

	MulticastMessage(manager.clients, ServerResponse)

	ServerResponse = PacketUPL2{
		Cmd: "gmsg",
		Val: GetGMSG(),
	}
	MulticastMessage(manager.clients, ServerResponse)

	ServerResponse = GetULIST(nil)
	MulticastMessage(manager.clients, ServerResponse)

	ServerResponse = PacketUPL2{
		Cmd:    "statuscode",
		Code:   "I:100 | OK",
		CodeID: 100,
	}
	MulticastMessage(manager.clients, ServerResponse)
}

func HandleGMSG(manager *Manager, msg interface{}) {
	log.Printf("Recived new client packet of type GMSG.")
	log.Printf("Message is %s", msg)
	MulticastMessage(manager.clients, AddGMSG(msg))
}
func HandleGVAR(manager *Manager, name interface{}, value interface{}) {
	log.Printf("Recived new client packet of type GVAR.")
	log.Printf("Var has a name of %s, and a value of %s.", name, value)
	MulticastMessage(manager.clients, AddGVAR(name, value))
}
func HandleSetID(manager *Manager, value interface{}, listener interface{}) {
	log.Printf("Recived new client packet of type SetID.")
	log.Printf("ID is '%s'.", value)
	if listener != nil {
		MulticastMessage(manager.clients, AddUser(value, listener))
	} else {
		MulticastMessage(manager.clients, AddUser(value, nil))
	}
	MulticastMessage(manager.clients, PacketUPL2{
		Cmd:    "statuscode",
		Code:   "I:100 | OK",
		CodeID: 100,
	})
}
func HandlePMSG(manager *Manager, value interface{}, clients map[*Client]Client) {
	if containsValueClientlist(manager.clients, clients) {
		MulticastMessage(manager.clients, PacketUPL2{
			Cmd:    "statuscode",
			Code:   "E:103 | IDNotFound	",
			CodeID: 100,
		})
	} else {
		MulticastMessage(manager.clients, PacketUPL2{
			Cmd:    "statuscode",
			Code:   "E:103 | IDNotFound	",
			CodeID: 100,
		})
	}
}
