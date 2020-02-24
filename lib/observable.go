// Package lib contains convenience functionality for repeated tasks, e.g. implementing the observer pattern.
package lib

import (
	"fmt"
	"sync"
)

type observers map[string][]chan bool

// ConcurrentObservable enables KeyObservable functionality.
// It may be used as a mixin to enable observability.
type ConcurrentObservable struct {
	DataLock     sync.RWMutex
	ObserverLock sync.Mutex
	observers
}

// Key may be any type that can be marshaled into plain-text
type Key fmt.Stringer

// PlainKey is a box for plain strings so that they may be used as Key type.
type PlainKey struct {
	Text string
}

// String marshals PlainKey into a string.
func (pk PlainKey) String() string {
	return pk.Text
}

// NewPlainKey constructs a PlainKey from a plain string.
func NewPlainKey(s string) PlainKey {
	return PlainKey{Text:s}
}

// KeyObservable is a collection of key-addressable things, where each key may be observed for changes.
type KeyObservable interface {
	// Observe registers a new observer for the given key.
	// The returned channel will receive a bool each time the observed key changes.
	Observe(k Key) chan bool
	// Notify notifies all registered observers of the given key of a change to the key.
	Notify(k Key)
}

func (l ConcurrentObservable) Observe(k Key) chan bool {
	l.ObserverLock.Lock()
	defer l.ObserverLock.Unlock()
	observer := make(chan bool)
	key := k.String()
	l.observers[key] = append(l.observers[key], observer)
	return observer
}

func (l ConcurrentObservable) Notify(k Key) {
	l.ObserverLock.Lock()
	defer l.ObserverLock.Unlock()
	for _, observer := range l.observers[k.String()] {
		// The channel may block if no-one is observing
		go (func() {
			observer <- true
		})()
	}
}

// NewConcurrentObservable constructs a new ConcurrentObservable.
func NewConcurrentObservable() ConcurrentObservable {
	return ConcurrentObservable{
		DataLock:     sync.RWMutex{},
		ObserverLock: sync.Mutex{},
		observers:    make(observers),
	}
}
