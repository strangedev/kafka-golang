package local

import (
	"encoding/json"
	"errors"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/google/uuid"
	"github.com/linkedin/goavro"
	"github.com/strangedev/kafka-golang/consumer/router"
	"github.com/strangedev/kafka-golang/schema"
	"github.com/strangedev/kafka-golang/schema/command"
	"github.com/strangedev/kafka-golang/schema/maps"
	"log"
)

type Repo struct {
	Schemata maps.SchemaMap
	Aliases  maps.AliasMap
	router.TopicRouter
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

	codec, err := goavro.NewCodec(request.Spec)
	if err != nil {
		log.Println(err)
		return err
	}

	repo.Schemata.Insert(request.UUID, codec)

	return nil
}

func (repo Repo) handleAliasUpdate(event *kafka.Message) error {
	var request command.AliasRequest
	err := json.Unmarshal(event.Value, &request)
	if err != nil {
		return err
	}

	log.Printf("^^ AliasRequest %v: %v\n", request.UUID, request.Alias)

	repo.Aliases.Insert(schema.Alias(request.Alias), request.UUID)

	return nil
}

func NewLocalRepo(consumer *kafka.Consumer) (Repo, error) {
	repo := Repo{
		TopicRouter: router.NewTopicRouter(consumer),
		Schemata:    maps.NewSchemaMap(),
		Aliases:     maps.NewAliasMap(),
	}
	log.Printf("Created schema repository with TopicRouter %v", repo.TopicRouter)

	repo.NewRoute("schema_update", repo.handleSchemaUpdate)
	repo.NewRoute("schema_alias", repo.handleAliasUpdate)

	return repo, nil
}
