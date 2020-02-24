package main

import (
	"flag"
	"github.com/strangedev/kafka-golang/schema"
	"github.com/strangedev/kafka-golang/schema/repo"
	"github.com/strangedev/kafka-golang/utils"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var broker string

func init() {
	flag.StringVar(&broker, "broker", "broker0:9092", "URL of a Kafka broker")
}

func main() {
	flag.Parse()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	log.Printf("Pulling schemata from %v", broker)

	schemaRepo, err := repo.NewLocalRepo(broker)
	utils.CheckFatal("Unable to initialize schema repository", err)

	var stop chan bool
	if len(os.Args) == 2 {
		schemaAlias := schema.Alias(os.Args[1])

		log.Printf("Wait for schema alias %v\n", schemaAlias)
		aliasReady := schemaRepo.WaitAliasReady(schemaAlias)

		stop, err := schemaRepo.Run()
		utils.CheckFatal("Unable to start schema repository", err)

		ok := <-utils.SigAbort(aliasReady, signals)
		if ok {
			log.Println("Alias ready")
		}
		stop <- true
		os.Exit(0)
	}
	stop, err = schemaRepo.Run()
	utils.CheckFatal("Unable to start schema repository", err)
	log.Print("Waiting for signal...")
	<-signals
	stop <- true
}
