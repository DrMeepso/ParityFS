package server

import (
	common "ParityFS/Common"
	"errors"

	"go.etcd.io/bbolt"
)

func (server *IServer) LoginWithCredentials(username, password string) (bool, string) {
	// We need to check if the creds are valid

	err := server.boltDB.View(func(tx *bbolt.Tx) error {
		// Get the credentials bucket
		bucket := tx.Bucket([]byte("Credentials"))
		if bucket == nil {
			return errors.New("credentials bucket not found")
		}

		// Check if the username exists
		hashedPassword := bucket.Get([]byte(username))
		if hashedPassword == nil {
			return errors.New("username does not exist")
		}

		// Verify the password
		if common.HashPassword(password) != string(hashedPassword) {
			return errors.New("invalid password")
		}

		return nil
	})

	if err != nil {
		Log("Error logging in user:", err)
		return false, err.Error()
	}

	Log("User logged in successfully:", username)
	return true, "User logged in successfully"

}

func (server *IServer) RegisterWithCredentials(username, password string) (bool, string) {
	// We need to check if the creds are valid

	err := server.boltDB.Update(func(tx *bbolt.Tx) error {

		// Get the credentials bucket
		bucket := tx.Bucket([]byte("Credentials"))
		if bucket == nil {
			return errors.New("credentials bucket not found")
		}

		// Check if the username already exists
		if bucket.Get([]byte(username)) != nil {
			return errors.New("username already exists")
		}

		// Store the hashed password
		hashedPassword := common.HashPassword(password)
		if err := bucket.Put([]byte(username), []byte(hashedPassword)); err != nil {
			return err
		}

		return nil

	})

	if err != nil {
		Log("Error registering user:", err)
		return false, err.Error()
	}

	Log("User registered successfully:", username)
	return true, "User registered successfully"

}

func (server *IServer) DoseUserExist(username string) (bool, string) {
	// Check if the user exists in the database
	var exists bool
	err := server.boltDB.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte("Credentials"))
		if bucket == nil {
			return errors.New("credentials bucket not found")
		}
		if bucket.Get([]byte(username)) != nil {
			exists = true
		} else {
			exists = false
		}
		return nil
	})

	if err != nil {
		Log("Error checking user existence:", err)
		return false, err.Error()
	}

	return exists, ""
}
