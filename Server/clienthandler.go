package server

import (
	common "ParityFS/Common"
	"net"
	"reflect"
)

type RemoteClient struct {
	Conn     net.Conn
	Address  string
	Username string

	ParentServer *IServer
}

func (Client *RemoteClient) SendPacket(data any) {
	// Send a packet to the client
	if err := common.SendPacket(Client.Conn, data); err != nil {
		common.RemoteLog("\033[91mServer >\033[0m", "Error sending packet to client:", err)
	}
}

func HandleNewClient(conn net.Conn, address string, Server IServer) *RemoteClient {
	client := &RemoteClient{
		Conn:         conn,
		Address:      address,
		Username:     "",
		ParentServer: &Server,
	}

	Server.addConnectedClient(client)

	go ClientLoop(client, Server)

	return client
}

func ClientLoop(client *RemoteClient, Server IServer) {
	// Handle client communication here
	// This is a placeholder for the actual implementation
	// You would typically read from client.Conn and respond accordingly
	// For now, we just log the new client connection
	Log("New client connected: ", client.Address)

	aknowledgment := common.JoinAknowledgment{
		ServerProtocolVersion: Server.Version,
	}
	if err := common.SendPacket(client.Conn, aknowledgment); err != nil {
		Log("Error sending acknowledgment to client:", err)
		return
	}

	defer func() {
		client.Conn.Close()
		Server.removeConnectedClient(client)
		Log("Client disconnected: ", client.Address)
	}()

	for {

		// Placeholder for reading from the client
		// In a real implementation, you would read data from client.Conn
		// and process it accordingly
		buffer := make([]byte, 1024)
		n, err := client.Conn.Read(buffer)
		if err != nil {
			Log("Error reading from client: ", err)
			return
		}

		if n > 0 {
			Packet, err := common.ReadPacket(buffer[:n])
			if err != nil {
				Log("Error reading packet from client: ", err)
				continue
			}

			switch data := Packet.Data.(type) {
			case *common.LoginRequest:
				{
					client.Username = data.Username
					succ, msg := Server.LoginWithCredentials(data.Username, data.Password)
					client.SendPacket(common.LoginResponse{Success: succ, Message: msg})

					if !succ {
						client.Conn.Close()
						Log("Client disconnected due to failed login: ", client.Address)
					} else {
						Log("User [", client.Username, "] logged in from ", client.Address)
					}

				}
			default:
				Log("Received unknown packet command from client: ", Packet.Command)
				Log("Packet data type:", reflect.TypeOf(Packet.Data))
			}
		}

	}

}
