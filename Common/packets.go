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

func SendPacket(socket net.Conn, data any) error {

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
			return err
		}
	} else {
		RemoteLog("\033[91mCommon >\033[0m", "Error: Packet type not registered:", dataType)
		return &net.OpError{
			Op:  "send",
			Net: "tcp",
			Err: &json.UnsupportedTypeError{Type: dataType},
		}
	}

	return nil
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

func ReadPacket(b []byte) (Packet, error) {
	var packet Packet

	// Unmarshal the JSON data into the Packet struct
	err := json.Unmarshal(b, &packet)
	if err != nil {
		RemoteLog("\033[91mCommon >\033[0m", "Error unmarshalling JSON:", err)
		return packet, err
	}

	// Check if the command is registered
	if _, ok := CommandToPacket[packet.Command]; !ok {
		RemoteLog("\033[91mCommon >\033[0m", "Error: Command not registered:", packet.Command)
		return packet, &json.UnmarshalTypeError{Value: packet.Command}
	}

	data, err := json.Marshal(packet.Data)
	if err != nil {
		RemoteLog("\033[91mCommon >\033[0m", "Error marshalling packet data:", err)
		return packet, err
	}

	// Unmarshal the data into the appropriate type
	packetType := CommandToPacket[packet.Command]
	packetValue := reflect.New(packetType).Interface()
	err = json.Unmarshal(data, packetValue)
	if err != nil {
		RemoteLog("\033[91mCommon >\033[0m", "Error unmarshalling packet data:", err)
		return packet, err
	}

	packet.Data = packetValue

	return packet, nil
}

func RegisterPackets() {
	PacketToCommand = make(map[reflect.Type]string)
	CommandToPacket = make(map[string]reflect.Type)

	// Register all packet types here
	AddPacketType("JoinAknowledgment", JoinAknowledgment{})
	AddPacketType("LoginRequest", LoginRequest{})
	AddPacketType("LoginResponse", LoginResponse{})

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
