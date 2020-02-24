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
	"github.com/strangedev/kafka-golang/consumer"
	"github.com/strangedev/kafka-golang/lib"
	"sync"
)

// Handler is a function used in a KeyRouter which handles a single kind of message.
// The Handler may return an error that will be logged.
type Handler func(event *kafka.Message) error

// KeyRouter is a generic router for Kafka which routes consumed messages to a set of configured handlers based on a key.
type KeyRouter interface {
	consumer.Consumer
	// NewRoute registers a handler for the given key. Only one handler may exist per key.
	NewRoute(key lib.Key, handler Handler)
	// Handle routes a single Event either to the appropriate Handler function or to a generic error handler.
	// Note that this method may block.
	Handle(event kafka.Event)
}

// ConcurrentRouter is a wrapper around a generic Router which controls concurrent access.
// This allows one to add new routes even while handling events. Simply wrap your existing Router
// in a ConcurrentRouter.
type ConcurrentRouter struct {
	Router KeyRouter
	sync.RWMutex
}

func (c ConcurrentRouter)NewRoute(key lib.Key, handler Handler) {
	c.Lock()
	defer c.Unlock()
	c.Router.NewRoute(key, handler)
}

func (c ConcurrentRouter) Handle(event kafka.Event) {
	c.RLock()
	defer c.RUnlock()
	c.Router.Handle(event)
}

func (c ConcurrentRouter)Run() (stop chan bool, err error) {
	stop, err = c.Router.Run()
	return
}