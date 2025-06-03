package server

import (
	common "ParityFS/Common"
	"crypto/tls"
	"net"
	"os"
	"strconv"

	"go.etcd.io/bbolt"
)

func Log(args ...any) {
	common.RemoteLog("\033[92mServer >\033[0m ", args...)
}

// should be paths to certificate files
type ICertificate struct {
	crt string
	key string
}

type IServer struct {
	port int
	host string

	// certificate for TLS
	certificate ICertificate

	Version     int
	tlslistener *net.Listener

	connectedClients []*RemoteClient

	BoltDB *bbolt.DB
}

var (
	Server IServer = IServer{

		port: 51888,
		host: "localhost",

		certificate: ICertificate{
			crt: "./server.crt",
			key: "./server.key",
		},

		Version:          common.ProtocallVersion,
		tlslistener:      nil,
		connectedClients: make([]*RemoteClient, 0, 100),

		BoltDB: nil,
	}
)

func ServerMain() {

	Log("ParityFS In Server Mode")

	// read the cmd line arguments for server options
	args := os.Args[1:]

	for i := 0; i < len(args); i++ {

		switch args[i] {

		case "--dev":
			fallthrough
		case "--allow-invalid-cert":
			fallthrough
		case "--server":
			break

		case "-p":
			fallthrough
		case "--port":
			i++
			if i < len(args) {
				err := error(nil)
				Server.port, err = strconv.Atoi(args[i])
				if err != nil {
					Log("Invalid port number:", args[i])
					return
				}
			} else {
				Log("No port specified, using default port 51888")
			}

		case "-h":
			fallthrough
		case "--host":
			i++
			if i < len(args) {
				Server.host = args[i]
			} else {
				Log("No host specified, using default host localhost")
			}

		case "--certificate":
			fallthrough
		case "--cert":
			fallthrough
		case "-c":
			i++
			if i < len(args) {
				Server.certificate.crt = args[i]
			} else {
				Log("No certificate file specified, using default empty certificate")
			}
		case "--key":
			fallthrough
		case "-k":
			i++
			if i < len(args) {
				Server.certificate.key = args[i]
			} else {
				Log("No key file specified, using default empty key")
			}

		default:
			Log("Unknown argument: ", args[i])
		}

	}

	Log("Server Configuration:")
	Log("  Port:", Server.port)
	Log("  Host:", Server.host)
	Log("  Certificate:", Server.certificate.crt)
	Log("  Key:", Server.certificate.key)
	Log("  Server Version:", Server.Version)

	// check if parityfs.db exists, if not create it
	if _, err := os.Stat("parityfs.db"); os.IsNotExist(err) {
		Log("Database file does not exist, creating parityfs.db")
		common.CreateDB()
	}

	Server.BoltDB = common.OpenDB()
	if Server.BoltDB == nil {
		Log("Error opening database, exiting server.")
		return
	}

	// check if certificate files exist
	if _, err := os.Stat(Server.certificate.crt); os.IsNotExist(err) {
		Log("Certificate file does not exist:", Server.certificate.crt)
		Log("Please generate a certificate using --certgen or provide a valid certificate file.")
		return
	}
	if _, err := os.Stat(Server.certificate.key); os.IsNotExist(err) {
		Log("Key file does not exist:", Server.certificate.key)
		Log("Please generate a certificate using --certgen or provide a valid key file.")
		return
	}

	// create a TLS instance
	cert, err := tls.LoadX509KeyPair(Server.certificate.crt, Server.certificate.key)
	if err != nil {
		Log("Error loading certificate and key:", err.Error())
		return
	}

	config := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}

	listener, err := tls.Listen("tcp", Server.host+":"+strconv.Itoa(Server.port), config)
	if err != nil {
		Log("Error starting server:", err.Error())
		return
	}

	Log("Server is listening on ", Server.host+":"+strconv.Itoa(Server.port))
	Server.tlslistener = &listener

	for {
		conn, err := (*Server.tlslistener).Accept()
		if err != nil {
			Log("Error accepting connection:", err.Error())
			continue
		}

		address := conn.RemoteAddr().String()
		HandleNewClient(conn, address, Server)
	}
}

func (server *IServer) Close() {
	if server.tlslistener != nil {
		(*server.tlslistener).Close()
		Log("Server listener closed")
	}

	if server.BoltDB != nil {
		server.BoltDB.Close()
		Log("Database connection closed")
	}

	for _, client := range server.connectedClients {
		client.Conn.Close()
		Log("Closed connection for client:", client.Address)
	}

	Log("Server shutdown complete")
}

func (server *IServer) addConnectedClient(client *RemoteClient) {
	// Add a client to the connected clients map
	if client != nil && client.Address != "" {
		server.connectedClients = append(server.connectedClients, client)
		Log("Client added:", client.Address)
	} else {
		Log("Invalid client, cannot add to connected clients")
	}
}
func (server *IServer) removeConnectedClient(client *RemoteClient) {
	// Remove a client from the connected clients map
	if client != nil && client.Address != "" {
		for i, c := range server.connectedClients {
			if c.Conn == client.Conn {
				server.connectedClients = append(server.connectedClients[:i], server.connectedClients[i+1:]...)
				Log("Client removed:", client.Address)
				return
			}
		}
		Log("Client not found in connected clients:", client.Address)
	} else {
		Log("Invalid client, cannot remove from connected clients")
	}
}
