/* Copyright 2020 Noah Hummel
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
 */

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
