package server

import (
	common "ParityFS/Common"
	"net"
)

type RemoteClient struct {
	Conn     net.Conn
	Address  string
	Username string
}

func HandleNewClient(conn net.Conn, address string, Server IServer) *RemoteClient {
	client := &RemoteClient{
		Conn:     conn,
		Address:  address,
		Username: "",
	}

	go ClientLoop(client, Server)

	return client
}

func ClientLoop(client *RemoteClient, Server IServer) {
	// Handle client communication here
	// This is a placeholder for the actual implementation
	// You would typically read from client.Conn and respond accordingly
	// For now, we just log the new client connection
	Log("New client connected:", client.Address)

	aknowledgment := common.JoinAknowledgment{
		ServerProtocolVersion: Server.Version,
	}
	if err := common.SendJSON(client.Conn, aknowledgment); err != nil {
		Log("Error sending acknowledgment to client:", err)
		client.Conn.Close()
		return
	}

	defer func() {
		client.Conn.Close()
		delete(Server.connectedClients, client.Address)
		Log("Client disconnected:", client.Address)
	}()

	for {

		// Placeholder for reading from the client
		// In a real implementation, you would read data from client.Conn
		// and process it accordingly
		buffer := make([]byte, 1024)
		n, err := client.Conn.Read(buffer)
		if err != nil {
			Log("Error reading from client:", err)
			return
		}

		if n > 0 {

		}

	}

}
