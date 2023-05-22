package main

func HandleHandshake(manager Manager) {

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
		Cmd: "direct",
		Val: "I:100 | OK",
	}
	MulticastMessage(manager.clients, ServerResponse)

	ServerResponse = PacketUPL2{
		Cmd: "direct",
		Val: "I:100 | OK",
	}
	MulticastMessage(mgr.clients, ServerResponse)
}
