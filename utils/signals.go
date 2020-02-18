package utils

import "os"

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
