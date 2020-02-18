package observable

import (
	"fmt"
	"sync"
)

type observers map[string][]chan bool

type LockableObservable struct {
	DataLock     sync.RWMutex
	ObserverLock sync.Mutex
	Observers    observers
}

type Key fmt.Stringer

type BasicObservable interface {
	Observe(k Key) chan bool
	Notify(k Key)
}

func (l LockableObservable) Observe(k Key) chan bool {
	l.ObserverLock.Lock()
	defer l.ObserverLock.Unlock()
	observer := make(chan bool)
	key := k.String()
	l.Observers[key] = append(l.Observers[key], observer)
	return observer
}

func (l LockableObservable) Notify(k Key) {
	l.ObserverLock.Lock()
	defer l.ObserverLock.Unlock()
	for _, observer := range l.Observers[k.String()] {
		observer <- true
	}
}

func NewLockableObservable() LockableObservable {
	return LockableObservable{
		DataLock:     sync.RWMutex{},
		ObserverLock: sync.Mutex{},
		Observers:    make(observers),
	}
}
