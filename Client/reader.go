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
			}

		default:
			Log("Received unknown packet command:", Packet.Command)
			Log("Packet data type:", reflect.TypeOf(Packet.Data))

		}

	}

}
