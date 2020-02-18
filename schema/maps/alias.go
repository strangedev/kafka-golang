package maps

import (
	"github.com/google/uuid"
	"kafka-go/lib/observable"
	"kafka-go/schema"
)

type aliasMapType map[schema.Alias]uuid.UUID

type AliasMap struct {
	observable.LockableObservable
	Map aliasMapType
}

func (m AliasMap) Insert(alias schema.Alias, schemaUUID uuid.UUID) bool {
	_, overwritten := m.Map[alias]
	m.Map[alias] = schemaUUID
	m.Notify(alias)
	return overwritten
}

func NewAliasMap() AliasMap {
	return AliasMap{
		LockableObservable: observable.NewLockableObservable(),
		Map:                make(aliasMapType),
	}
}