package main

import (
	"kafka-go/schema"
	"kafka-go/schema/query/local"
	"kafka-go/schema/query/versioned"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	repo, err := local.NewRepo("broker0:9092")
	if err != nil {
		log.Fatalln("Can't initialize Repo")
	}
	versions := versioned.FromRepo(repo)
	schemaVersion := schema.NameVersion{Name: "mySchema", Version: 13}
	schemaReady := versions.WaitSchemaReady(schemaVersion)
	stop := repo.Run()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	go (func() {
		<-signals
		stop <- true
	})()

	ok := <-schemaReady
	if ok {
		log.Println("Schema ready")
		uuid, _ := versions.WhoIs(schema.Alias(schemaVersion.String()))
		log.Printf("Is UUID %v\n", uuid)
	} else {
		log.Println("Aborted")
	}
}
