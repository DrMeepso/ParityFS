package main

import (
	client "ParityFS/Client"
	common "ParityFS/Common"
	server "ParityFS/Server"
	"os"
	"time"
)

func main() {

	common.RegisterPackets() // Register all packet types

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

		case "--dev":
			println("Running in development mode")
			common.IsDevelopmentMode = true

			defer func() {

				// If the program panics, print the error
				if r := recover(); r != nil {
					println("Error:", r.(string))
				}

				println("Somthing went wrong, dumping logs...")

				if len(common.BufferedLogs) > 0 {
					for _, log := range common.BufferedLogs {
						println(log)
					}
					common.BufferedLogs = make([]string, 0, 100) // clear the buffer
				}

			}()

			go server.ServerMain()
			time.Sleep(100 * time.Millisecond) // Give server time to start
			go client.ClientMain()
			common.HandelLogging()
			return

		}
	}

	client.ClientMain()

}
