package main

import (
	"fmt"
	"github.com/google/uuid"
	"kafka-go/schema"
	"os"
)

func main() {
	cmd, err := schema.NewCommander("broker0:9092")
	if err != nil {
		fmt.Println("Can't create commander")
		fmt.Println(err)
		os.Exit(1)
	}
	if os.Args[1] == "spec" {
		schemaUUID, err := uuid.Parse(os.Args[2])
		if err != nil {
			fmt.Println("Can't parse UUID")
			fmt.Println(err)
			os.Exit(1)
		}
		err = cmd.UpdateSchema(schemaUUID, os.Args[3])
		if err != nil {
			fmt.Println("Can't send command")
			fmt.Println(err)
			os.Exit(1)
		}
	} else if os.Args[1] == "alias" {
		schemaUUID, err := uuid.Parse(os.Args[3])
		if err != nil {
			fmt.Println("Can't parse UUID")
			fmt.Println(err)
			os.Exit(1)
		}
		err = cmd.UpdateAlias(os.Args[2], schemaUUID)
		if err != nil {
			fmt.Println("Can't send command")
			fmt.Println(err)
			os.Exit(1)
		}
	}

}
