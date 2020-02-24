package utils

import "os"

// SigAbort waits for a value on the ready channel.
// Waiting is be aborted, when a value is received on the abort channel first.
// It returns a channel that receives a value when either condition has been met,
// true if ready was first to receive, true otherwise.
func SigAbort(ready chan bool, abort chan os.Signal) chan bool {
	fin := make(chan bool)
	go (func() {
		select {
		case _ = <-abort:
			fin <- false
			close(fin)
			return
		case _ = <-ready:
			fin <- true
			close(fin)
			return
		}
	})()
	return fin
}
