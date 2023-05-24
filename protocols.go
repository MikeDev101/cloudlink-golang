package main

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

	ServerResponse = GetULIST()
	MulticastMessage(manager.clients, ServerResponse)

	ServerResponse = PacketUPL2{
		Cmd:    "statuscode",
		Code:   "I:100 | OK",
		CodeID: 100,
	}
	MulticastMessage(manager.clients, ServerResponse)
}

func HandleGMSG(manager *Manager, msg interface{}) {
	MulticastMessage(manager.clients, AddGMSG(msg))
}
func HandleGVAR(manager *Manager, name interface{}, value interface{}) {
	MulticastMessage(manager.clients, AddGVAR(name, value))
}
func HandleSetID(manager *Manager, value interface{}) {
	MulticastMessage(manager.clients, AddUser(value))
	MulticastMessage(manager.clients, PacketUPL2{
		Cmd:    "statuscode",
		Code:   "I:100 | OK",
		CodeID: 100,
	})
}
