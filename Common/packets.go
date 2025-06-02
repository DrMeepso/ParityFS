package common

import (
	"encoding/json"
	"net"
	"reflect"
)

var (
	CommandToPacket map[string]reflect.Type = map[string]reflect.Type{}
	PacketToCommand map[reflect.Type]string = map[reflect.Type]string{}
)

func AddPacketType(Command string, Type any) {
	CommandToPacket[Command] = reflect.TypeOf(Type)
	PacketToCommand[reflect.TypeOf(Type)] = Command
}

func RegisterPackets() {
	// Register all packet types here
	AddPacketType("JoinAknowledgment", JoinAknowledgment{})
	AddPacketType("LoginRequest", LoginRequest{})
	AddPacketType("LoginResponse", LoginResponse{})
}

func SendPacket(socket net.Conn, data any) {

	//get the type of the data
	dataType := reflect.TypeOf(data)

	//check if the type is registered
	if command, ok := PacketToCommand[dataType]; ok {
		// create a packet with the command and data
		packet := Packet{
			Command: command,
			Data:    data,
		}

		// send the packet as JSON
		if err := SendJSON(socket, packet); err != nil {
			RemoteLog("\033[91mCommon >\033[0m", "Error sending packet:", err)
		}
	} else {
		RemoteLog("\033[91mCommon >\033[0m", "Error: Packet type not registered:", dataType)
	}

}

func SendJSON(socket net.Conn, data any) error {

	jsonData, err := json.Marshal(data)
	if err != nil {
		RemoteLog("\033[91mCommon >\033[0m", "Error marshalling JSON:", err)
	}

	_, err = socket.Write(jsonData)
	if err != nil {
		RemoteLog("\033[91mCommon >\033[0m", "Error writing to socket:", err)
		return err
	}

	return nil

}

// Packet represents a generic packet structure used for communication
// All packets use this structure as a base
// The `Command` field specifies the type of packet being sent
// The `Data` field contains the payload of the packet, which can be any type
type Packet struct {
	Command string `json:"command"`
	Data    any    `json:"data,omitempty"`
}

// Server > Client
// Sent when a client joins the server
type JoinAknowledgment struct {
	ServerProtocolVersion int `json:"server_protocol_version"`
}

// Client > Server
// Sent after the aknowledgment packet
// Contains the client's protocol version and credentials
type LoginRequest struct {
	Username              string `json:"username"`
	Password              string `json:"password"`
	ClientProtocolVersion int    `json:"client_protocol_version"`
}

// Server > Client
// Sent in response to a login request
type LoginResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}
