package schema

import (
	"encoding/json"
	"errors"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/google/uuid"
	"github.com/linkedin/goavro"
	"log"
	"os"
)

type Repo interface {
	Decode(schema uuid.UUID, datum []byte) (interface{}, error)
	Encode(schema uuid.UUID, datum interface{}) ([]byte, error)
	WaitSchemaReady(schema uuid.UUID) chan bool
	WaitAliasReady(alias string) chan bool
	ListSchemata() []uuid.UUID
	ListAliases() []string
	WhoIs(alias string) (uuid.UUID, bool)
	GetSpecification(schema uuid.UUID) (specification string, ok bool)
	Count() int
}

type localRepo struct {
	Schemata        map[uuid.UUID]*goavro.Codec
	SchemaObservers map[uuid.UUID][]chan bool
	Alias           map[string]uuid.UUID
	AliasObservers  map[string][]chan bool
	Consumer        *kafka.Consumer
	SchemaCount     int
}

func (repo *localRepo) schemaUpdate(schemaUUID uuid.UUID, schemaSpecification string) error {
	codec, err := goavro.NewCodec(schemaSpecification)
	if err != nil {
		log.Println(err)
		return err
	}
	_, exists := repo.Schemata[schemaUUID]
	if !exists {
		repo.SchemaCount++
	}
	repo.Schemata[schemaUUID] = codec
	return nil
}

func (repo *localRepo) aliasUpdate(schemaUUID uuid.UUID, alias string) {
	repo.Alias[alias] = schemaUUID
}

func (repo localRepo) Decode(schema uuid.UUID, datum []byte) (interface{}, error) {
	codec, ok := repo.Schemata[schema]
	if !ok {
		return nil, errors.New("schema not present")
	}
	native, _, err := codec.NativeFromBinary(datum)
	return native, err
}

func (repo localRepo) Encode(schema uuid.UUID, datum interface{}) ([]byte, error) {
	codec, ok := repo.Schemata[schema]
	if !ok {
		return nil, errors.New("schema not present")
	}
	binary, err := codec.BinaryFromNative(nil, datum)
	return binary, err
}

func (repo localRepo) WaitSchemaReady(schema uuid.UUID) chan bool {
	schemaIsReady := make(chan bool)
	_, ok := repo.Schemata[schema]
	if ok {
		log.Println("-> schema is already present")
		go (func() {
			schemaIsReady <- true
		})()
	} else {
		log.Println("-> establish SchemaObserver")
		repo.SchemaObservers[schema] = append(repo.SchemaObservers[schema], schemaIsReady)
	}
	return schemaIsReady
}

func (repo localRepo) WaitAliasReady(alias string) chan bool {
	aliasIsReady := make(chan bool)
	schemaUUID, ok := repo.Alias[alias]
	if ok {
		log.Println("-> alias is already present")
		go (func() {
			log.Println("-> alias is already present")
			schemaIsReady := repo.WaitSchemaReady(schemaUUID)
			<-schemaIsReady
			aliasIsReady <- true
		})()
	} else {
		log.Println("-> establish AliasObserver")
		repo.AliasObservers[alias] = append(repo.AliasObservers[alias], aliasIsReady)
	}
	return aliasIsReady
}

func (repo localRepo) ListSchemata() []uuid.UUID {
	schemata := make([]uuid.UUID, 0, len(repo.Schemata))
	for schemaUUID := range repo.Schemata {
		schemata = append(schemata, schemaUUID)
	}
	return schemata
}

func (repo localRepo) ListAliases() []string {
	aliases := make([]string, 0, len(repo.Alias))
	for alias := range repo.Alias {
		aliases = append(aliases, alias)
	}
	return aliases
}

func (repo localRepo) WhoIs(alias string) (uuid.UUID, bool) {
	value, ok := repo.Alias[alias]
	return value, ok
}

func (repo localRepo) GetSpecification(schema uuid.UUID) (specification string, ok bool) {
	var codec *goavro.Codec
	codec, ok = repo.Schemata[schema]
	if !ok {
		return "", false
	}
	return codec.Schema(), true
}

func (repo localRepo) Count() int {
	return repo.SchemaCount
}

type UpdateRequest struct {
	UUID uuid.UUID `json:"UUID"`
	Spec string    `json:"spec"`
}

type AliasRequest struct {
	UUID  uuid.UUID `json:"UUID"`
	Alias string    `json:"alias"`
}

func (repo localRepo) handleSchemaUpdate(event *kafka.Message) error {
	var request UpdateRequest
	err := json.Unmarshal(event.Value, &request)
	if err != nil {
		return err
	}
	log.Printf("^^ UpdateRequest %v: %v\n", request.UUID, request.Spec)
	err = repo.schemaUpdate(request.UUID, request.Spec)
	if err != nil {
		return err
	}

	log.Println("-> Notify AliasObservers")
	for _, observer := range repo.SchemaObservers[request.UUID] {
		observer <- true
	}

	return nil
}

func (repo localRepo) handleAliasUpdate(event *kafka.Message) error {
	var request AliasRequest
	err := json.Unmarshal(event.Value, &request)
	if err != nil {
		return err
	}

	log.Printf("^^ AliasRequest %v: %v\n", request.UUID, request.Alias)
	repo.aliasUpdate(request.UUID, request.Alias)

	log.Println("-> Notify AliasObservers")
	for _, observer := range repo.AliasObservers[request.Alias] {
		log.Println("--> Check alias present")
		schemaUUID, ok := repo.Alias[request.Alias]
		if ok {
			log.Println("--> Alias present, check schema present")
			go (func() {
				schemaIsReady := repo.WaitSchemaReady(schemaUUID)
				<-schemaIsReady
				log.Println("--> Schema is ready, notify AliasObserver")
				observer <- true
			})()
		}
	}

	return nil
}

func NewRepo(broker string, signals chan os.Signal) (Repo, error) {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":     broker,
		"group.id":              uuid.New().String(),
		"broker.address.family": "v4",
		"session.timeout.ms":    6000,
		"auto.offset.reset":     "earliest",
	})
	if err != nil {
		log.Println(err)
		return nil, err
	}

	repo := localRepo{
		Schemata:        make(map[uuid.UUID]*goavro.Codec),
		SchemaObservers: make(map[uuid.UUID][]chan bool, 1),
		Alias:           make(map[string]uuid.UUID),
		AliasObservers:  make(map[string][]chan bool, 1),
		Consumer:        consumer,
		SchemaCount:     0,
	}

	topics := []string{"schema_update", "schema_alias"}
	err = consumer.SubscribeTopics(topics, nil)
	go (func() {
		run := true
		for run {
			select {
			case sig := <-signals:
				log.Println(sig)
				run = false
			default:
				event := consumer.Poll(100)
				if event == nil {
					continue
				}

				switch e := event.(type) {
				case *kafka.Message:
					switch *e.TopicPartition.Topic {
					case "schema_update":
						err = repo.handleSchemaUpdate(e)
						if err != nil {
							log.Printf("!! Error handling SchemaUpdate: %v", err)
						}

					case "schema_alias":
						err = repo.handleAliasUpdate(e)
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
	return repo, nil
}
