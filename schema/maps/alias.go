package maps

import (
	"github.com/google/uuid"
	"github.com/strangedev/kafka-golang/lib"
	"github.com/strangedev/kafka-golang/schema"
)

type aliasMapType map[schema.Alias]uuid.UUID

type AliasMap struct {
	lib.LockableObservable
	Map aliasMapType
}

func (m AliasMap) Insert(alias schema.Alias, schemaUUID uuid.UUID) bool {
	m.DataLock.Lock()
	_, overwritten := m.Map[alias]
	m.Map[alias] = schemaUUID
	m.DataLock.Unlock()
	m.Notify(alias)
	return overwritten
}

func NewAliasMap() AliasMap {
	return AliasMap{
		LockableObservable: lib.NewLockableObservable(),
		Map:                make(aliasMapType),
	}
}
