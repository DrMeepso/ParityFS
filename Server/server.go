package server

import (
	"ParityFS/common"
	"crypto/tls"
	"net"
	"os"
	"strconv"
)

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

	serverVersion int
	tlslistener   *net.Listener
}

var (
	Server IServer = IServer{

		port: 51888,
		host: "localhost",

		certificate: ICertificate{
			crt: "./server.crt",
			key: "./server.key",
		},

		serverVersion: common.ProtocallVersion,
	}
)

func ServerMain() {

	println("ParityFS In Server Mode")

	// read the cmd line arguments for server options
	args := os.Args[1:]

	for i := 0; i < len(args); i++ {

		switch args[i] {

		case "--server":
			break

		case "-p":
		case "--port":
			i++
			if i < len(args) {
				err := error(nil)
				Server.port, err = strconv.Atoi(args[i])
				if err != nil {
					println("Invalid port number:", args[i])
					return
				}
			} else {
				println("No port specified, using default port 51888")
			}

		case "--host":
			i++
			if i < len(args) {
				Server.host = args[i]
			} else {
				println("No host specified, using default host localhost")
			}

		case "--certificate":
		case "--cert":
		case "-c":
			i++
			if i < len(args) {
				Server.certificate.crt = args[i]
			} else {
				println("No certificate file specified, using default empty certificate")
			}
		case "--key":
		case "-k":
			i++
			if i < len(args) {
				Server.certificate.key = args[i]
			} else {
				println("No key file specified, using default empty key")
			}

		default:
			println("Unknown argument:", args[i])
		}

	}

	println("Server Configuration:")
	println("  Port:", Server.port)
	println("  Host:", Server.host)
	println("  Certificate:", Server.certificate.crt)
	println("  Key:", Server.certificate.key)
	println("  Server Version:", Server.serverVersion)

	// check if certificate files exist
	if _, err := os.Stat(Server.certificate.crt); os.IsNotExist(err) {
		println("Certificate file does not exist:", Server.certificate.crt)
		println("Please generate a certificate using --certgen or provide a valid certificate file.")
		return
	}
	if _, err := os.Stat(Server.certificate.key); os.IsNotExist(err) {
		println("Key file does not exist:", Server.certificate.key)
		println("Please generate a certificate using --certgen or provide a valid key file.")
		return
	}

	// create a TLS instance
	cert, err := tls.LoadX509KeyPair(Server.certificate.crt, Server.certificate.key)
	if err != nil {
		println("Error loading certificate and key:", err.Error())
		return
	}

	config := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}

	listener, err := tls.Listen("tcp", Server.host+":"+strconv.Itoa(Server.port), config)
	if err != nil {
		println("Error starting server:", err.Error())
		return
	}

	println("Server is listening on", Server.host+":"+strconv.Itoa(Server.port))
	Server.tlslistener = &listener

}
