package producer

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

// KafkaProducer is a wrapper for kafka.Producer which handles delivery failure in a canonical way.
type KafkaProducer struct {
	Producer *kafka.Producer
}

// NewKafkaProducer constructs a new KafkaProducer that will produce into the given Kafka broker.
func NewKafkaProducer(broker string) (*KafkaProducer, error) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": broker})
	if err != nil {
		return nil, err
	}
	return &KafkaProducer{Producer: p}, nil
}

// LowLevelProducer provides access to synchronous sending routines.
type LowLevelProducer interface {
	// ProduceSimpleSync sends a message without headers or key into a specific topic and partition.
	// It is mainly used for testing.
	ProduceSimpleSync(topic string, partition int32, value []byte) error
	// ProduceSync synchronously produces a Message.
	ProduceSync(message *kafka.Message) error
}

func (k *KafkaProducer) ProduceSync(message *kafka.Message) error {
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
