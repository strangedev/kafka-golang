package main

import (
	"encoding/json"
	"flag"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/google/uuid"
	"github.com/strangedev/kafka-golang/schema"
	"github.com/strangedev/kafka-golang/schema/query"
	"github.com/strangedev/kafka-golang/utils"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

type SchemaDTO struct {
	UUID          uuid.UUID `json:"uuid"`
	Specification string    `json:"spec"`
}

type SchemataDTO struct {
	Schemata []SchemaDTO `json:"schemata"`
}

type SchemaListDTO struct {
	Count    int         `json:"count"`
	Schemata []uuid.UUID `json:"schemata"`
}

type AliasDTO struct {
	Alias schema.Alias	`json:"alias"`
	UUID  uuid.UUID `json:"uuid"`
}

type AliasListDTO struct {
	Aliases []schema.Alias `json:"aliases"`
	Count   int      `json:"count"`
}

type AliasesDTO struct {
	Aliases []AliasDTO `json:"aliases"`
}

var broker string

func init() {
	flag.StringVar(&broker, "broker", "broker0:9092", "URL of a Kafka broker")
}

func writeJSON(writer http.ResponseWriter, data interface{}) {
	ret, err := json.Marshal(data)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		log.Println(err)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	writer.Header().Set("Access-Control-Allow-Headers", "*")
	_, err = writer.Write(ret)
}

func main() {
	flag.Parse()

	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":     broker,
		"group.id":              uuid.New().String(),
		"broker.address.family": "v4",
		"session.timeout.ms":    6000,
		"auto.offset.reset":     "earliest",
	})
	utils.CheckFatal("Unable to initialize Kafka consumer", err)

	schemaRepo, err := query.NewLocalRepo(consumer)
	utils.CheckFatal("Unable to initialize schema repository", err)

	stop, err := schemaRepo.Run()
	utils.CheckFatal("Unable to start schema repository", err)

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	go (func() {
		for sig := range signals {
			stop<-true
			log.Panicf("Caught %v", sig)
		}
	})()

	http.HandleFunc("/schema/list", func(writer http.ResponseWriter, request *http.Request) {
		schemata := schemaRepo.ListSchemata()
		schemaList := SchemaListDTO{Schemata: schemata, Count: len(schemata)}

		writeJSON(writer, schemaList)
	})

	http.HandleFunc("/schema/describe", func(writer http.ResponseWriter, request *http.Request) {
		params := request.URL.Query()
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			log.Println(err)
			return
		}

		schemaUUIDs := params["uuid"]
		if len(schemaUUIDs) < 1 {
			http.Error(writer, "Required params <uuid>", http.StatusBadRequest)
			return
		}

		schemata := SchemataDTO{Schemata: make([]SchemaDTO, 0, len(schemaUUIDs))}
		for _, uuidString := range schemaUUIDs {
			schemaUUID, err := uuid.Parse(uuidString)
			if err != nil {
				http.Error(writer, err.Error(), http.StatusBadRequest)
				log.Println(err)
				return
			}

			spec, ok := schemaRepo.GetSpecification(schemaUUID)
			if !ok {
				log.Println(err)
				continue
			}

			schemata.Schemata = append(schemata.Schemata, SchemaDTO{UUID: schemaUUID, Specification: spec})
		}

		writeJSON(writer, schemata)
	})

	http.HandleFunc("/alias/list", func(writer http.ResponseWriter, request *http.Request) {
		aliases := schemaRepo.ListAliases()
		aliasList := AliasListDTO{Aliases: aliases, Count: len(aliases)}

		writeJSON(writer, aliasList)
	})

	http.HandleFunc("/alias/describe", func(writer http.ResponseWriter, request *http.Request) {
		params := request.URL.Query()
		if err != nil {
			log.Println(err)
			return
		}

		aliasesQuery := params["alias"]
		if len(aliasesQuery) < 1 {
			http.Error(writer, "Required params <alias>", http.StatusBadRequest)
			return
		}

		aliases := AliasesDTO{Aliases: make([]AliasDTO, 0)}
		for _, aliasString := range aliasesQuery {
			alias := schema.Alias(aliasString)
			schemaUUID, ok := schemaRepo.WhoIs(alias)
			if !ok {
				log.Println(err)
				continue
			}

			aliases.Aliases = append(aliases.Aliases, AliasDTO{UUID: schemaUUID, Alias: alias})
		}

		writeJSON(writer, aliases)
	})

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Println(err)
		return
	}
}
