package main

import (
	"fmt"
	"kafka-golang/schema"
	"kafka-golang/schema/query/local"
	"kafka-golang/utils"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	schemaRepo, err := local.NewRepo("broker0:9092")
	if err != nil {
		log.Fatalln("Can't initialize SchemaRepo")
	}

	stop := schemaRepo.Run()

	if len(os.Args) == 2 {
		schemaAlias := schema.Alias(os.Args[1])

		fmt.Printf("Wait for schema alias %v\n", schemaAlias)
		aliasReady := schemaRepo.WaitAliasReady(schemaAlias)

		ok := <-utils.SigAbort(aliasReady, signals)
		if ok {
			fmt.Println("Alias ready")
		}
		stop <- true
		os.Exit(0)
	}
	<-signals
	stop <- true
}
