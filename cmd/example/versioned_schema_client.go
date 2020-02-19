package main

import (
	"flag"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/google/uuid"
	"github.com/strangedev/kafka-golang/schema"
	"github.com/strangedev/kafka-golang/schema/query/local"
	"github.com/strangedev/kafka-golang/schema/query/versioned"
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

	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":     broker,
		"group.id":              uuid.New().String(),
		"broker.address.family": "v4",
		"session.timeout.ms":    6000,
		"auto.offset.reset":     "earliest",
	})
	utils.CheckFatal("Unable to initialize Kafka consumer", err)

	schemaRepo, err := local.NewLocalRepo(consumer)
	utils.CheckFatal("Unable to initialize schema repository", err)

	versions := versioned.FromRepo(schemaRepo)
	schemaVersion := schema.NameVersion{Name: "mySchema", Version: 13}
	schemaReady := versions.WaitSchemaReady(schemaVersion)
	stop, err := schemaRepo.Run()
	utils.CheckFatal("Unable to start schema repository", err)

	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	go (func() {
		<-signals
		stop <- true
	})()

	ok := <-schemaReady
	if ok {
		log.Println("Schema ready")
		schemaUUID, _ := versions.WhoIs(schema.Alias(schemaVersion.String()))
		log.Printf("Is UUID %v\n", schemaUUID)
	} else {
		log.Println("Aborted")
	}
}
