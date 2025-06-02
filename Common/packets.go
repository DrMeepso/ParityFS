package common

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
