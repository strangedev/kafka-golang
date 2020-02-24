package main

import (
	"encoding/json"
	"flag"
	"github.com/google/uuid"
	"github.com/strangedev/kafka-golang/schema"
	"github.com/strangedev/kafka-golang/schema/explorer"
	"github.com/strangedev/kafka-golang/schema/repo"
	"github.com/strangedev/kafka-golang/utils"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

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
	log.Printf("Broker: %v", broker)
	
	schemaRepo, err := repo.NewLocalRepo(broker)
	utils.CheckFatal("Unable to initialize schema repository", err)

	stop, err := schemaRepo.Run()
	utils.CheckFatal("Unable to start schema repository", err)

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	go (func() {
		for sig := range signals {
			stop <- true
			log.Panicf("Caught %v", sig)
		}
	})()

	http.HandleFunc("/schema/list", func(writer http.ResponseWriter, request *http.Request) {
		schemata := schemaRepo.ListSchemata()
		schemaList := explorer.SchemaListDTO{Schemata: schemata, Count: len(schemata)}

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

		schemata := explorer.SchemataDTO{Schemata: make([]explorer.SchemaDTO, 0, len(schemaUUIDs))}
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

			schemata.Schemata = append(schemata.Schemata, explorer.SchemaDTO{UUID: schemaUUID, Specification: spec})
		}

		writeJSON(writer, schemata)
	})

	http.HandleFunc("/alias/list", func(writer http.ResponseWriter, request *http.Request) {
		aliases := schemaRepo.ListAliases()
		aliasList := explorer.AliasListDTO{Aliases: aliases, Count: len(aliases)}

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

		aliases := explorer.AliasesDTO{Aliases: make([]explorer.AliasDTO, 0)}
		for _, aliasString := range aliasesQuery {
			alias := schema.Alias(aliasString)
			schemaUUID, ok := schemaRepo.WhoIs(alias)
			if !ok {
				log.Println(err)
				continue
			}

			aliases.Aliases = append(aliases.Aliases, explorer.AliasDTO{UUID: schemaUUID, Alias: alias})
		}

		writeJSON(writer, aliases)
	})

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Println(err)
		return
	}
}
