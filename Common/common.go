package common

import (
	"fmt"
	"sync"
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
