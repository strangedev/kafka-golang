package utils

import "log"

func CheckFatal(message string, err error) {
	if err != nil {
		log.Fatalf("!! %v Error: %v", message, err)
	}
}
