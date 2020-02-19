package command

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/google/uuid"
	"github.com/strangedev/kafka-golang/producer"
)

type Commander struct {
	*producer.KafkaProducer
}

type Updater interface {
	UpdateSchema(schemaUUID uuid.UUID, specification string) error
	UpdateAlias(alias string, schemaUUID uuid.UUID) error
}

func NewUpdater(broker string) (Updater, error) {
	p, err := producer.NewKafkaProducer(broker)
	if err != nil {
		return nil, err
	}
	return NewUpdaterWithProducer(p), nil
}

func NewUpdaterWithProducer(p *producer.KafkaProducer) Updater {
	return Commander{p}
}

func (cmd Commander) UpdateSchema(schemaUUID uuid.UUID, specification string) error {
	topic := "schema_update"
	request := UpdateRequest{UUID: schemaUUID, Spec: specification}
	return cmd.ProduceJSONSync(topic, kafka.PartitionAny, request)
}

func (cmd Commander) UpdateAlias(alias string, schemaUUID uuid.UUID) error {
	topic := "schema_alias"
	request := AliasRequest{UUID: schemaUUID, Alias: alias}
	return cmd.ProduceJSONSync(topic, kafka.PartitionAny, request)
}
