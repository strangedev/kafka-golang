package local

import (
	"encoding/json"
	"errors"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/google/uuid"
	"github.com/linkedin/goavro"
	"kafka-go/schema"
	"kafka-go/schema/command"
	"kafka-go/schema/maps"
	"log"
)

type Repo struct {
	Schemata maps.SchemaMap
	Aliases  maps.AliasMap
	Consumer *kafka.Consumer
}

func (repo *Repo) schemaUpdate(schemaUUID uuid.UUID, schemaSpecification string) error {
	codec, err := goavro.NewCodec(schemaSpecification)
	if err != nil {
		log.Println(err)
		return err
	}
	repo.Schemata.Insert(schemaUUID, codec)
	return nil
}

func (repo *Repo) aliasUpdate(schemaUUID uuid.UUID, alias schema.Alias) {
	repo.Aliases.Insert(alias, schemaUUID)
}

func (repo Repo) Decode(schema uuid.UUID, datum []byte) (interface{}, error) {
	codec, ok := repo.Schemata.Map[schema]
	if !ok {
		return nil, errors.New("schema not present")
	}
	native, _, err := codec.NativeFromBinary(datum)
	return native, err
}

func (repo Repo) Encode(schema uuid.UUID, datum interface{}) ([]byte, error) {
	codec, ok := repo.Schemata.Map[schema]
	if !ok {
		return nil, errors.New("schema not present")
	}
	binary, err := codec.BinaryFromNative(nil, datum)
	return binary, err
}

func (repo Repo) WaitSchemaReady(schema uuid.UUID) chan bool {
	return repo.Schemata.Observe(schema)
}

func (repo Repo) WaitAliasReady(alias schema.Alias) chan bool {
	aliasIsReady := make(chan bool)
	go (func() {
		<-repo.Aliases.Observe(alias)
		schemaUUID, _ := repo.Aliases.Map[alias]
		schemaReady := repo.Schemata.Observe(schemaUUID)
		_, ok := repo.GetSpecification(schemaUUID)
		if !ok {
			<-schemaReady
		}
		aliasIsReady <- true
	})()
	return aliasIsReady
}

func (repo Repo) ListSchemata() []uuid.UUID {
	schemata := make([]uuid.UUID, 0, len(repo.Schemata.Map))
	for schemaUUID := range repo.Schemata.Map {
		schemata = append(schemata, schemaUUID)
	}
	return schemata
}

func (repo Repo) ListAliases() []schema.Alias {
	aliases := make([]schema.Alias, 0, len(repo.Aliases.Map))
	for alias := range repo.Aliases.Map {
		aliases = append(aliases, alias)
	}
	return aliases
}

func (repo Repo) WhoIs(alias schema.Alias) (uuid.UUID, bool) {
	value, ok := repo.Aliases.Map[alias]
	return value, ok
}

func (repo Repo) GetSpecification(schema uuid.UUID) (string, bool) {
	codec, ok := repo.Schemata.Map[schema]
	if !ok {
		return "", false
	}
	return codec.Schema(), true
}

func (repo Repo) Count() int {
	return len(repo.Schemata.Map)
}

func (repo Repo) handleSchemaUpdate(event *kafka.Message) error {
	var request command.UpdateRequest
	err := json.Unmarshal(event.Value, &request)
	if err != nil {
		return err
	}
	log.Printf("^^ UpdateRequest %v: %v\n", request.UUID, request.Spec)
	err = repo.schemaUpdate(request.UUID, request.Spec)
	if err != nil {
		return err
	}

	return nil
}

func (repo Repo) handleAliasUpdate(event *kafka.Message) error {
	var request command.AliasRequest
	err := json.Unmarshal(event.Value, &request)
	if err != nil {
		return err
	}

	log.Printf("^^ AliasRequest %v: %v\n", request.UUID, request.Alias)
	repo.aliasUpdate(request.UUID, schema.Alias(request.Alias))

	return nil
}

func (repo Repo) Run() (stop chan bool) {
	stop = make(chan bool, 1)
	go (func() {
		run := true
		for run {
			select {
			case sig := <-stop:
				log.Println(sig)
				run = false
			default:
				event := repo.Consumer.Poll(100)
				if event == nil {
					continue
				}

				switch e := event.(type) {
				case *kafka.Message:
					switch *e.TopicPartition.Topic {
					case "schema_update":
						err := repo.handleSchemaUpdate(e)
						if err != nil {
							log.Printf("!! Error handling SchemaUpdate: %v", err)
						}
					case "schema_alias":
						err := repo.handleAliasUpdate(e)
						if err != nil {
							log.Printf("!! Error handling AliasUpdate: %v", err)
						}
					}
				case *kafka.Error:
					log.Printf("!! Kafka Error: %v: %v\n", e.Code(), e)
				default:
					log.Printf("-- Ignored %v\n", e)
				}
			}
		}
	})()
	return stop
}

func NewRepo(broker string) (Repo, error) {
	repo := Repo{
		Schemata: maps.NewSchemaMap(),
		Aliases:  maps.NewAliasMap(),
		Consumer: nil,
	}
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":     broker,
		"group.id":              uuid.New().String(),
		"broker.address.family": "v4",
		"session.timeout.ms":    6000,
		"auto.offset.reset":     "earliest",
	})
	if err != nil {
		log.Println(err)
		return repo, err
	}
	repo.Consumer = consumer

	topics := []string{"schema_update", "schema_alias"}
	err = consumer.SubscribeTopics(topics, nil)
	return repo, nil
}
