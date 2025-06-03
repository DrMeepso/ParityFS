package main

import (
	client "ParityFS/Client"
	common "ParityFS/Common"
	server "ParityFS/Server"
	"os"
	"time"

	"go.etcd.io/bbolt"
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

		case "--adduser":
			if i+2 >= len(args) {
				println("Usage: --adduser <username> <password>")
				return
			}
			username, password := args[i+1], args[i+2]
			i += 2 // Skip the next two arguments

			fakeServer := server.IServer{
				BoltDB: common.OpenDB(),
			}

			success, message := fakeServer.CreateUser(username, password)
			if success {
				println("User created successfully:", username)
			} else {
				println("Error creating user:", message)
			}

		case "--listusers":
			fakeServer := server.IServer{
				BoltDB: common.OpenDB(),
			}
			err := fakeServer.BoltDB.View(func(tx *bbolt.Tx) error {
				bucket := tx.Bucket([]byte("Credentials"))
				if bucket == nil {
					println("No users found.")
					return nil
				}

				println("Registered users:")
				bucket.ForEach(func(k, v []byte) error {
					println(string(k)) // Print the username
					return nil
				})
				return nil
			})

			if err != nil {
				println("Error listing users:", err.Error())
			}

		case "--removeuser":
			if i+1 >= len(args) {
				println("Usage: --removeuser <username>")
				return
			}
			username := args[i+1]
			i++ // Skip the next argument

			fakeServer := server.IServer{
				BoltDB: common.OpenDB(),
			}
			err := fakeServer.BoltDB.View(func(tx *bbolt.Tx) error {
				bucket := tx.Bucket([]byte("Credentials"))
				if bucket == nil {
					return nil // No users to remove
				}

				if bucket.Get([]byte(username)) == nil {
					println("User not found:", username)
					return nil
				}

				if err := bucket.Delete([]byte(username)); err != nil {
					return err
				}

				println("User removed successfully:", username)
				return nil
			})

			if err != nil {
				println("Error removing user:", err.Error())
			}

		}
	}

	client.ClientMain()

}
