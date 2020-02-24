package utils

import "log"

// CheckFatal checks the given error, logs it in a canonical format and exists the program.
func CheckFatal(message string, err error) {
	if err != nil {
		log.Fatalf("!! %v Error: %v", message, err)
	}
}
