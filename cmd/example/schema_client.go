package main

import (
	"fmt"
	"kafka-go/schema"
	"kafka-go/utils"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	schemaAlias := os.Args[1]
	schemaRepo, err := schema.NewRepo("broker0:9092", signals)
	if err != nil {
		log.Fatalln("Can't initialize SchemaRepo")
	}

	fmt.Printf("Wait for schema alias %v\n", schemaAlias)
	aliasReady := schemaRepo.WaitAliasReady(schemaAlias)

	<-utils.SigAbort(aliasReady, signals)
	fmt.Println("Alias ready")


}
