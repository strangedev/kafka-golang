package main

import (
	"encoding/json"
	"github.com/google/uuid"
	"kafka-go/schema/query/local"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

// TODO

type Schema struct {
	UUID          uuid.UUID `json:"uuid"`
	Specification string    `json:"spec"`
}

type Schemata struct {
	Schemata []Schema `json:"schemata"`
}

type SchemaList struct {
	Count    int         `json:"count"`
	Schemata []uuid.UUID `json:"schemata"`
}

type Alias struct {
	Alias string    `json:"alias"`
	UUID  uuid.UUID `json:"uuid"`
}

type AliasList struct {
	Aliases []string `json:"aliases"`
	Count   int      `json:"count"`
}

type Aliases struct {
	Aliases []Alias `json:"aliases"`
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
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	go (func() {
		for sig := range signals {
			log.Panicf("Caught %v", sig)
		}
	})()

	schemaRepo, err := local.NewRepo("broker0:9092", signals)
	if err != nil {
		log.Fatalln("Can't initialize SchemaRepo")
	}

	http.HandleFunc("/schema/list", func(writer http.ResponseWriter, request *http.Request) {
		schemata := schemaRepo.ListSchemata()
		schemaList := SchemaList{Schemata: schemata, Count: len(schemata)}

		writeJSON(writer, schemaList)
	})

	http.HandleFunc("/schema/describe", func(writer http.ResponseWriter, request *http.Request) {
		query := request.URL.Query()
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			log.Println(err)
			return
		}

		schemaUUIDs := query["uuid"]
		if len(schemaUUIDs) < 1 {
			http.Error(writer, "Required query <uuid>", http.StatusBadRequest)
			return
		}

		schemata := Schemata{Schemata: make([]Schema, 0, len(schemaUUIDs))}
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

			schemata.Schemata = append(schemata.Schemata, Schema{UUID: schemaUUID, Specification: spec})
		}

		writeJSON(writer, schemata)
	})

	http.HandleFunc("/alias/list", func(writer http.ResponseWriter, request *http.Request) {
		aliases := schemaRepo.ListAliases()
		aliasList := AliasList{Aliases: aliases, Count: len(aliases)}

		writeJSON(writer, aliasList)
	})

	http.HandleFunc("/alias/describe", func(writer http.ResponseWriter, request *http.Request) {
		query := request.URL.Query()
		if err != nil {
			log.Println(err)
			return
		}

		aliasesQuery := query["alias"]
		if len(aliasesQuery) < 1 {
			http.Error(writer, "Required query <alias>", http.StatusBadRequest)
			return
		}

		aliases := Aliases{Aliases: make([]Alias, 0)}
		for _, alias := range aliasesQuery {
			schemaUUID, ok := schemaRepo.WhoIs(alias)
			if !ok {
				log.Println(err)
				continue
			}

			aliases.Aliases = append(aliases.Aliases, Alias{UUID: schemaUUID, Alias: alias})
		}

		writeJSON(writer, aliases)
	})

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Println(err)
		return
	}
}
