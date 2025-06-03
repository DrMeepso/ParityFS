package client

import (
	common "ParityFS/Common"
	"reflect"
)

func BeginReading() {

	for {

		// read from the server connection
		buffer := make([]byte, 1024)
		n, err := ServerConn.Read(buffer)
		if err != nil {
			Log("Error reading from server:", err)
			return
		}

		if n == 0 {
			Log("Server connection closed")
			return
		}

		Packet, err := common.ReadPacket(buffer[:n])
		if err != nil {
			Log("Error reading packet:", err)
			continue
		}

		switch data := Packet.Data.(type) {
		case *common.JoinAknowledgment:
			{
				Log("Received JoinAknowledgment from server")
				Log("Server Protocol Version:", data.ServerProtocolVersion)
				if data.ServerProtocolVersion != common.ProtocallVersion {
					Log("Protocol version mismatch! Client version:", common.ProtocallVersion, "Server version:", data.ServerProtocolVersion)
					Log("Please update your client or server to match the protocol version.")
					return
				}

				loginAttempt := common.LoginRequest{
					Username: Client.Credential.Username,
					Password: Client.Credential.Password,
				}

				if err := common.SendPacket(ServerConn, loginAttempt); err != nil {
					Log("Error sending login request:", err)
					return
				}

			}
		case *common.LoginResponse:
			{
				Log("Received LoginResponse from server")
				if data.Success {
					Log("Login successful!")
					Client.IsLoggedIn = true
				} else {
					Log("Login failed:", data.Message)
					return
				}
			}

		default:
			Log("Received unknown packet command:", Packet.Command)
			Log("Packet data type:", reflect.TypeOf(Packet.Data))

		}

	}

}
