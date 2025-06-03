package client

import (
	common "ParityFS/Common"
	"crypto/tls"
	"os"
	"reflect"
	"strconv"
)

func Log(args ...any) {
	common.RemoteLog("\033[36mClient >\033[0m ", args...)
}

type IServerInfo struct {
	Host string
	Port int
}

type ICredential struct {
	Username string
	Password string
}

type IClient struct {
	ServerInfo       IServerInfo
	Version          int
	Credential       ICredential
	IsLoggedIn       bool
	allowInvalidCert bool // for testing purposes, should be false in production
}

var Client IClient = IClient{
	ServerInfo: IServerInfo{
		Host: "localhost",
		Port: 51888,
	},
	Version: common.ProtocallVersion,
	Credential: ICredential{
		Username: "user",
		Password: "password",
	},
	IsLoggedIn: false,

	allowInvalidCert: false, // this is for testing purposes, should be false in production
}

var ServerConn *tls.Conn

func ClientMain() {
	Log("ParityFS In Client Mode")

	args := os.Args[1:]
	for i := 0; i < len(args); i++ {

		switch args[i] {

		case "--dev":
			break

		case "-user":
			fallthrough
		case "--username":
			i++
			if i < len(args) {
				Client.Credential.Username = args[i]
			} else {
				Log("Error: No username provided")
				return
			}
			break

		case "-pass":
			fallthrough
		case "--password":
			i++
			if i < len(args) {
				Client.Credential.Password = args[i]
			} else {
				Log("Error: No password provided")
				return
			}
			break

		case "-h":
			fallthrough
		case "--host":
			i++
			if i < len(args) {
				Client.ServerInfo.Host = args[i]
			} else {
				Log("Error: No host provided")
				return
			}
			break
		case "-p":
			fallthrough
		case "--port":
			i++
			if i < len(args) {
				var err error
				Client.ServerInfo.Port, err = strconv.Atoi(args[i])
				if err != nil {
					Log("Error: Invalid port number", args[i])
					return
				}
			} else {
				Log("Error: No port provided")
				return
			}
			break

		case "--allow-invalid-cert":
			Client.allowInvalidCert = true
			Log("Warning: Allowing invalid server certificate for testing purposes. This should not be used in production.")

		default:
			Log("Unknown argument: ", args[i])
			break

		}
	}

	Log("Client Configuration:")
	Log("  Host:", Client.ServerInfo.Host)
	Log("  Port:", Client.ServerInfo.Port)
	Log("  Username:", Client.Credential.Username)
	Log("  Password:", Client.Credential.Password)
	Log("  Version:", Client.Version)

	tlsConfig := &tls.Config{
		InsecureSkipVerify: Client.allowInvalidCert, // I dont own the server certificate, so I skip verification
	}

	conn, err := tls.Dial("tcp", Client.ServerInfo.Host+":"+strconv.Itoa(Client.ServerInfo.Port), tlsConfig)
	if err != nil {
		Log("Error connecting to server:", err.Error())
		if reflect.TypeOf(err).String() == "*tls.CertificateVerificationError" {
			Log("This is likely due to an invalid or self-signed certificate. Use --allow-invalid-cert to skip verification.")
			Log("If you have not had this issue before, please check your server's certificate and ensure it is valid.")
		} else {
			Log("Please check your server address and port.")
		}
		panic("Failed to connect to server: " + err.Error())
		return
	}

	ServerConn = conn

	defer conn.Close()
	Log("Connected to server: ", Client.ServerInfo.Host, ":", Client.ServerInfo.Port)

	// should be a goroutine
	// but i havent added the FUSEFS yet, so I just use a blocking call
	BeginReading()

}
