package common

import (
	"crypto"
	"fmt"
	"sync"

	"go.etcd.io/bbolt"
)

const (
	ProtocallVersion = 1
)

var (
	IsDevelopmentMode = false
)

var (
	BufferedLogs = make([]string, 0, 100)
	LogLock      sync.Mutex
)

func RemoteLog(from string, args ...any) {
	// lock the log buffer
	LogLock.Lock()
	BufferedLogs = append(BufferedLogs, from+fmt.Sprint(args...))
	// unlock the log buffer
	LogLock.Unlock()
}

func log(args ...any) {
	RemoteLog("\033[91mCommon >\033[0m ", args...)
}

func HandelLogging() {
	for {
		LogLock.Lock()
		if len(BufferedLogs) > 0 {
			for _, log := range BufferedLogs {
				fmt.Println(log)
			}
			BufferedLogs = make([]string, 0, 100) // clear the buffer
		}
		LogLock.Unlock()
	}
}

func HashPassword(password string) string {
	// This is a placeholder for the actual password hashing logic
	hash := crypto.MD5.New()
	hash.Write([]byte(password))
	return fmt.Sprintf("%x", hash.Sum(nil))
}

func CreateDB() {
	// Create a new BoltDB database
	db, err := bbolt.Open("parityfs.db", 0600, nil)
	if err != nil {
		log("Error opening database:", err)
		return
	}

	// Create a bucket for storing user credentials
	db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("Credentials"))
		if err != nil {
			log("Error creating bucket:", err)
		}
		return nil
	})

	defer db.Close()
}

func OpenDB() *bbolt.DB {
	// create a bbolt.DB instance
	db, err := bbolt.Open("parityfs.db", 0600, nil)
	if err != nil {
		log("Error opening database:", err)
		return nil
	}
	return db
}
