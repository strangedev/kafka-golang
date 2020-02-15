package schema

import (
	"encoding/json"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/google/uuid"
)

type commander struct {
	Producer *kafka.Producer
}

type Commander interface {
	UpdateSchema(schemaUUID uuid.UUID, specification string) error
	UpdateAlias(alias string, schemaUUID uuid.UUID) error
}

func NewCommander(broker string) (Commander, error) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": broker})
	if err != nil {
		return nil, err
	}
	commander := commander{Producer: p}
	return commander, nil
}

func (cmd commander) UpdateSchema(schemaUUID uuid.UUID, specification string) error {
	topic := "schema_update"
	request := UpdateRequest{UUID: schemaUUID, Spec:specification}
	value, err := json.Marshal(request)
	if err != nil {
		return err
	}
	deliveryChan := make(chan kafka.Event)
	err = cmd.Producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          value,
	}, deliveryChan)

	e := <-deliveryChan
	m := e.(*kafka.Message)

	if m.TopicPartition.Error != nil {
		fmt.Printf("Delivery failed: %v\n", m.TopicPartition.Error)
	} else {
		fmt.Printf("Delivered message to topic %s [%d] at offset %v\n",
			*m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
	}

	close(deliveryChan)
	return err
}

func (cmd commander) UpdateAlias(alias string, schemaUUID uuid.UUID) error {
	topic := "schema_alias"
	request := AliasRequest{UUID: schemaUUID, Alias:alias}
	value, err := json.Marshal(request)
	if err != nil {
		return err
	}
	deliveryChan := make(chan kafka.Event)
	err = cmd.Producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          value,
	}, deliveryChan)

	e := <-deliveryChan
	m := e.(*kafka.Message)

	if m.TopicPartition.Error != nil {
		fmt.Printf("Delivery failed: %v\n", m.TopicPartition.Error)
	} else {
		fmt.Printf("Delivered message to topic %s [%d] at offset %v\n",
			*m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
	}

	close(deliveryChan)
	return err
}
