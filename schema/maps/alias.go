package maps

import (
	"github.com/google/uuid"
	"github.com/strangedev/kafka-golang/lib"
	"github.com/strangedev/kafka-golang/schema"
)

type aliasMapType map[schema.Alias]uuid.UUID

// SchemaMap is a KeyObservable map of Aliases to UUIDs
type AliasMap struct {
	lib.ConcurrentObservable
	Map aliasMapType
}

// Upsert inserts or updates an Alias, UUID pair into the map.
// All of the map's observers are notified of this change.
// It returns true, if the map entry did already exist and was overwritten.
func (m AliasMap) Insert(alias schema.Alias, schemaUUID uuid.UUID) bool {
	m.DataLock.Lock()
	_, overwritten := m.Map[alias]
	m.Map[alias] = schemaUUID
	m.DataLock.Unlock()
	m.Notify(alias)
	return overwritten
}

// NewAliasMap constructs an empty AliasMap with no observers.
func NewAliasMap() AliasMap {
	return AliasMap{
		ConcurrentObservable: lib.NewConcurrentObservable(),
		Map:                  make(aliasMapType),
	}
}
