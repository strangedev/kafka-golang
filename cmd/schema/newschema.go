package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/google/uuid"
	"github.com/linkedin/goavro"
	"github.com/strangedev/kafka-golang/schema"
	"github.com/strangedev/kafka-golang/schema/command"
	"github.com/strangedev/kafka-golang/schema/explorer"
	"github.com/strangedev/kafka-golang/utils"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var name, brokerURL, explorerURL string
var skipExplorerCheck bool

func init() {
	flag.StringVar(&name, "name", "", "A name for the new schema")
	flag.StringVar(&brokerURL, "broker", "broker0:9092", "URL of a Kafka broker")
	flag.StringVar(&explorerURL, "explorer", "schema-explorer:8085", "Use the schema explorer to check if the schema already exists before creating it.")
	flag.BoolVar(&skipExplorerCheck, "skip-check", false, "Do not use the schema explorer to check if if the schema already exists.")
}

func latestSchemaVersion() (uint, bool) {
	route := fmt.Sprintf("http://%v/alias/list", explorerURL)
	resp, err := http.Get(route)
	utils.CheckFatal("Unable to list current aliases from explorer", err)

	body, err := ioutil.ReadAll(resp.Body)
	var aliases explorer.AliasListDTO
	err = json.Unmarshal(body, &aliases)
	utils.CheckFatal("Unable to list current aliases from explorer (error unmarshalling response)", err)

	for _, alias := range aliases.Aliases {
		schemaVersion, err := schema.VersionFromAlias(alias)
		if err != nil {
			log.Printf("The alias %v does not seems to be in a versioned format. %v", alias, err)
			continue
		}
		if schemaVersion.Name == name {
			return schemaVersion.Version, true
		}
	}
	return 0, false
}

func main() {
	flag.Parse()
	schemaVersion := schema.NewVersionOrigin(name)
	if !skipExplorerCheck {
		log.Println("Checking if schema already exists...")
		if latestVersion, exists := latestSchemaVersion(); exists {
			log.Fatalf("A schema with that name already exists, its latest version is %v", latestVersion)
		}
		log.Println("Schema does not yet exist, continuing")
	}
	scanner := bufio.NewScanner(os.Stdin)
	buffer := bytes.Buffer{}
	for scanner.Scan() {
		buffer.Write(scanner.Bytes())
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("Unable to read from stdin: %v", err)
	}
	_, err := goavro.NewCodec(buffer.String())
	utils.CheckFatal("This does not seem like a valid Avro schema", err)

	schemaUUID := uuid.New()
	cmd, err := command.NewUpdater(brokerURL)
	utils.CheckFatal("Unable to initialize updater", err)
	err = cmd.UpdateSchema(schemaUUID, buffer.String())
	utils.CheckFatal("Unable to produce SchemaUpdate event", err)
	log.Printf("Updated schema %v", schemaUUID)
	err = cmd.UpdateAlias(schemaVersion.String(), schemaUUID)
	utils.CheckFatal("Unable to produce AliasUpdate event", err)
	log.Printf("Updated alias %v", schemaVersion.String())

	log.Println("The schema has been created.")
}
