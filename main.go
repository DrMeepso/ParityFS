package main

import (
	"ParityFS/common"
	"ParityFS/server"
	"os"
)

func main() {

	// get command line arguments
	args := os.Args[1:]

	for i := 0; i < len(args); i++ {

		switch args[i] {

		case "--server":
			server.ServerMain()
			return

		case "--certgen":
			// Generate a new certificate
			_, certPEM, keyPEM, err := common.CreateX509Pair()
			if err != nil {
				println("Error generating certificate:", err.Error())
				return
			}

			// Save the certificate and key to files
			err = os.WriteFile("server.crt", certPEM, 0644)
			if err != nil {
				println("Error writing certificate file:", err.Error())
				return
			}
			err = os.WriteFile("server.key", keyPEM, 0644)
			if err != nil {
				println("Error writing key file:", err.Error())
				return
			}
			println("New certificate and key generated: server.crt and server.key")
			return
		}

	}

	println("Assuming Client Mode...")

}
