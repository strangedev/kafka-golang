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

	schemaVersion := schema.NameVersion{Name: "mySchema", Version: 13}
	schemaReady := schemaRepo.WaitVersionReady(schemaVersion)
	stop, err := schemaRepo.Run()
	utils.CheckFatal("Unable to start schema repository", err)

	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	go (func() {
		<-signals
		stop <- true
	})()

	ok := <-schemaReady
	if ok {
		schemaUUID, _ := schemaRepo.WhoIs(schemaVersion.Alias())
		log.Printf("Schema ready, has UUID %v\n", schemaUUID)
	} else {
		log.Println("Aborted")
	}
}
