/* Copyright 2020 Noah Hummel
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
 */

package router

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/strangedev/kafka-golang/lib"
	"log"
)

// TopicRouter is a Router that routes based on the Kafka topic of consumed messages.
type TopicRouter struct {
	*kafka.Consumer
	Handlers map[lib.Key]Handler
}

func (t TopicRouter) NewRoute(topic lib.Key, handler Handler) {
	log.Printf("New handler for topic %v", topic)
	t.Handlers[topic] = handler
}

func (t TopicRouter) Handle(event kafka.Event) {
	switch e := event.(type) {
	case *kafka.Message:
		topic := lib.NewPlainKey(*e.TopicPartition.Topic)
		err := t.Handlers[topic](e)
		if err != nil {
			log.Printf("!! Error in route handler: %v", err.Error())
		}
	case *kafka.Error:
		log.Printf("!! Kafka Error: %v: %v\n", e.Code(), e)
	default:
		log.Printf("-- Ignored %v\n", e)
	}
}

func (t TopicRouter) Topics() []string {
	topics := make([]string, 0, len(t.Handlers))
	for key, _ := range t.Handlers {
		topics = append(topics, key.String())
	}
	return topics
}

func (t TopicRouter) Run() (chan bool, error) {
	log.Printf("Starting TopicRouter with %v and %v handlers", t.Consumer, len(t.Topics()))
	err := t.Consumer.SubscribeTopics(t.Topics(), nil)
	if err != nil {
		return nil, err
	}
	log.Printf("Subscribed to topics %v with consumer %v", t.Topics(), t.Consumer)

	stop := make(chan bool, 1)
	go (func() {
	forever:
		for {
			select {
			case <-stop:
				break forever
			default:
				event := t.Consumer.Poll(100)
				if event == nil {
					continue
				}
				go t.Handle(event)
			}
		}
	})()

	return stop, nil
}

func NewTopicRouter(consumer *kafka.Consumer) TopicRouter {
	log.Printf("Created new TopicRouter with consumer %v", consumer)
	return TopicRouter{
		Consumer: consumer,
		Handlers: make(map[lib.Key]Handler),
	}
}
