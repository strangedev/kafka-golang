// Package repo contains implementations of the various schema.*Repo interfaces.
package repo

import (
	"encoding/json"
	"errors"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/google/uuid"
	"github.com/linkedin/goavro"
	"github.com/strangedev/kafka-golang/consumer/router"
	"github.com/strangedev/kafka-golang/lib"
	"github.com/strangedev/kafka-golang/schema"
	"github.com/strangedev/kafka-golang/schema/command"
	"github.com/strangedev/kafka-golang/schema/maps"
	"log"
)

// LocalRepo is a local Kafka consumer that implements the various schema.*Repo interfaces.
type LocalRepo struct {
	Schemata maps.SchemaMap
	Aliases  maps.AliasMap
	router.TopicRouter
}

func (repo LocalRepo) Decode(schema uuid.UUID, datum []byte) (interface{}, error) {
	codec, ok := repo.Schemata.Map[schema]
	if !ok {
		return nil, errors.New("schema not present")
	}
	native, _, err := codec.NativeFromBinary(datum)
	return native, err
}

func (repo LocalRepo) Encode(schema uuid.UUID, datum interface{}) ([]byte, error) {
	codec, ok := repo.Schemata.Map[schema]
	if !ok {
		return nil, errors.New("schema not present")
	}
	binary, err := codec.BinaryFromNative(nil, datum)
	return binary, err
}

func (repo LocalRepo) WaitSchemaReady(schema uuid.UUID) chan bool {
	_, ok := repo.GetSpecification(schema)
	if !ok {
		return repo.Schemata.Observe(schema)
	}
	ready := make(chan bool)
	go (func() {
		ready <- true
	})()
	return ready
}

func (repo LocalRepo) WaitAliasReady(alias schema.Alias) chan bool {
	schemaUUID, ok := repo.WhoIs(alias)
	if !ok {
		aliasIsReady := make(chan bool)
		go (func() {
			<-repo.Aliases.Observe(alias)
			schemaUUID, _ := repo.Aliases.Map[alias]
			<-repo.WaitSchemaReady(schemaUUID)
			aliasIsReady <- true
		})()
		return aliasIsReady
	}

	return repo.WaitSchemaReady(schemaUUID)
}

func (repo LocalRepo) ListSchemata() []uuid.UUID {
	schemata := make([]uuid.UUID, 0, len(repo.Schemata.Map))
	for schemaUUID := range repo.Schemata.Map {
		schemata = append(schemata, schemaUUID)
	}
	return schemata
}

func (repo LocalRepo) ListAliases() []schema.Alias {
	aliases := make([]schema.Alias, 0, len(repo.Aliases.Map))
	for alias := range repo.Aliases.Map {
		aliases = append(aliases, alias)
	}
	return aliases
}

func (repo LocalRepo) WhoIs(alias schema.Alias) (uuid.UUID, bool) {
	value, ok := repo.Aliases.Map[alias]
	return value, ok
}

func (repo LocalRepo) GetSpecification(schema uuid.UUID) (string, bool) {
	codec, ok := repo.Schemata.Map[schema]
	if !ok {
		return "", false
	}
	return codec.Schema(), true
}

func (repo LocalRepo) Count() int {
	return len(repo.Schemata.Map)
}

func (repo LocalRepo) handleSchemaUpdate(message *kafka.Message) error {
	var request command.UpdateRequest
	err := json.Unmarshal(message.Value, &request)
	if err != nil {
		return err
	}

	log.Printf("^^ UpdateRequest %v: %v\n", request.UUID, request.Spec)

	codec, err := goavro.NewCodec(request.Spec)
	if err != nil {
		log.Println(err)
		return err
	}

	repo.Schemata.Upsert(request.UUID, codec)

	return nil
}

func (repo LocalRepo) handleAliasUpdate(message *kafka.Message) error {
	var request command.AliasRequest
	err := json.Unmarshal(message.Value, &request)
	if err != nil {
		return err
	}

	log.Printf("^^ AliasRequest %v: %v\n", request.UUID, request.Alias)

	repo.Aliases.Insert(schema.Alias(request.Alias), request.UUID)

	return nil
}

// NewLocalRepo constructs a LocalRepo configured for the specified Kafka broker.
// Note that since the repo is a Consumer, it needs to be started with Run() before it starts consuming.
func NewLocalRepo(broker string) (LocalRepo, error) {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":     broker,
		"group.id":              uuid.New().String(),
		"broker.address.family": "v4",
		"session.timeout.ms":    6000,
		"auto.offset.reset":     "earliest",
	})
	if err != nil {
		return LocalRepo{}, err
	}

	repo := LocalRepo{
		TopicRouter: router.NewTopicRouter(consumer),
		Schemata:    maps.NewSchemaMap(),
		Aliases:     maps.NewAliasMap(),
	}
	log.Printf("Created schema repository with TopicRouter %v", repo.TopicRouter)

	repo.NewRoute(lib.NewPlainKey("schema_update"), repo.handleSchemaUpdate)
	repo.NewRoute(lib.NewPlainKey("schema_alias"), repo.handleAliasUpdate)

	return repo, nil
}

func (repo LocalRepo) DecodeVersion(schema schema.NameVersion, datum []byte) (interface{}, error) {
	if schemaUUID, ok := repo.WhoIs(schema.Alias()); !ok {
		decoded, err := repo.Decode(schemaUUID, datum)
		return decoded, err
	}
	return nil, errors.New("schema not know to this repo")
}

func (repo LocalRepo) EncodeVersion(schema schema.NameVersion, datum interface{}) ([]byte, error) {
	if schemaUUID, ok := repo.WhoIs(schema.Alias()); ok {
		encoded, err := repo.Encode(schemaUUID, datum)
		return encoded, err
	}
	return nil, errors.New("schema not know to this repo")
}

func (repo LocalRepo) WaitVersionReady(schema schema.NameVersion) chan bool {
	return repo.WaitAliasReady(schema.Alias())
}
