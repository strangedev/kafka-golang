package router

import "github.com/confluentinc/confluent-kafka-go/kafka"

type Handler func (event *kafka.Message) error

type Router interface {
	NewRoute(key string, handler Handler)
	Handle(event kafka.Event)
}