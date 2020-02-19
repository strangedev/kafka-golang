package producer

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type KafkaProducer struct {
	Producer *kafka.Producer
}

func NewKafkaProducer(broker string) (*KafkaProducer, error) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": broker})
	if err != nil {
		return nil, err
	}
	return &KafkaProducer{Producer: p}, nil
}

type LowLevelProducer interface {
	ProduceSimpleSync(topic string, partition int32, value []byte) error
	ProduceSync(message *kafka.Message) error
}

func (k* KafkaProducer) ProduceSync(message *kafka.Message) error {
	deliveryChan := make(chan kafka.Event)
	defer close(deliveryChan)

	err := k.Producer.Produce(message, deliveryChan)

	e := <-deliveryChan
	m := e.(*kafka.Message)
	if m.TopicPartition.Error != nil {
		fmt.Printf("Delivery failed: %v\n", m.TopicPartition.Error)
	} else {
		fmt.Printf("Delivered message to topic %s [%d] at offset %v\n",
			*m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
	}

	return err
}

func (k *KafkaProducer) ProduceSimpleSync(topic string, partition int32, value []byte) error {
	deliveryChan := make(chan kafka.Event)
	defer close(deliveryChan)

	return k.ProduceSync(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: partition},
		Value:          value,
	})
}
