package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"kafka-go/schema"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

type Schema struct {
	UUID          uuid.UUID `json:"UUID"`
	Specification string    `json:"spec"`
}

type Schemata struct {
	Schemata []Schema `json:"schemata"`
}

type SchemaList struct {
	Count    int      	 `json:"count"`
	Schemata []uuid.UUID `json:"schemata"`
}

type Alias struct {
	Alias string `json:"alias"`
	UUID uuid.UUID `json:"uuid"`
}

type AliasList struct {
	Aliases []string `json:"aliases"`
	Count int `json:"count"`
}

type Aliases struct {
	Aliases []Alias `json:"aliases"`
}

type Error struct {
	Message string `json:"message"`
}

func main() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	go (func() {
		for sig := range signals {
			log.Fatalf("Caught %v", sig)
		}
	})()

	schemaRepo, err := schema.NewRepo("broker0:9092", signals)
	if err != nil {
		log.Fatalln("Can't initialize SchemaRepo")
	}

	http.HandleFunc("/schema/list", func(writer http.ResponseWriter, request *http.Request) {
		schemata := schemaRepo.ListSchemata()
		schemaList := SchemaList{Schemata: schemata, Count: len(schemata)}
		ret, err := json.Marshal(schemaList)
		if err != nil {
			log.Println(err)
			return
		}

		_, err = fmt.Fprintf(writer, string(ret))
		if err != nil {
			log.Println(err)
			return
		}
	})

	http.HandleFunc("/schema/describe", func(writer http.ResponseWriter, request *http.Request) {
		query := request.URL.Query()
		if err != nil {
			log.Println(err)
			return
		}

		schemaUUIDs := query["uuid"]
		if len(schemaUUIDs) < 1 {
			errorMessage := Error{Message:"No UUIDs received"}
			ret, err := json.Marshal(errorMessage)
			if err != nil {
				log.Println(err)
				return
			}
			_, err = fmt.Fprintf(writer, string(ret))
			if err != nil {
				log.Println(err)
			}
			return
		}
		schemata := Schemata{Schemata: make([]Schema, 0, len(schemaUUIDs))}
		for _, uuidString := range schemaUUIDs {
			schemaUUID, err := uuid.Parse(uuidString)
			if err != nil {
			 	log.Println(err)
			 	continue
			}

			spec, ok := schemaRepo.GetSpecification(schemaUUID)
			if !ok {
				log.Println(err)
				continue
			}
			schemata.Schemata = append(schemata.Schemata, Schema{UUID:schemaUUID, Specification:spec})
		}

		ret, err := json.Marshal(schemata)
		if err != nil {
			log.Println(err)
			return
		}

		_, err = fmt.Fprintf(writer, string(ret))
		if err != nil {
			log.Println(err)
			return
		}
	})

	http.HandleFunc("/alias/list", func(writer http.ResponseWriter, request *http.Request) {
		aliases := schemaRepo.ListAliases()
		aliasList := AliasList{Aliases: aliases, Count: len(aliases)}
		ret, err := json.Marshal(aliasList)
		if err != nil {
			log.Println(err)
			return
		}

		_, err = fmt.Fprintf(writer, string(ret))
		if err != nil {
			log.Println(err)
			return
		}
	})

	http.HandleFunc("/alias/describe", func(writer http.ResponseWriter, request *http.Request) {
		query := request.URL.Query()
		if err != nil {
			log.Println(err)
			return
		}

		aliasesQuery := query["alias"]
		if len(aliasesQuery) < 1 {
			errorMessage := Error{Message:"No aliases received"}
			ret, err := json.Marshal(errorMessage)
			if err != nil {
				log.Println(err)
				return
			}
			_, err = fmt.Fprintf(writer, string(ret))
			if err != nil {
				log.Println(err)
			}
			return
		}
		aliases := Aliases{Aliases: make([]Alias, 0, )}
		for _, alias := range aliasesQuery {

			schemaUUID, ok := schemaRepo.WhoIs(alias)
			if !ok {
				log.Println(err)
				continue
			}
			aliases.Aliases = append(aliases.Aliases, Alias{UUID:schemaUUID, Alias:alias})
		}

		ret, err := json.Marshal(aliases)
		if err != nil {
			log.Println(err)
			return
		}

		_, err = fmt.Fprintf(writer, string(ret))
		if err != nil {
			log.Println(err)
			return
		}
	})

	http.Handle("/", http.FileServer(http.Dir("./cmd/schema/ui/explorer")))

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Println(err)
		return
	}
}
