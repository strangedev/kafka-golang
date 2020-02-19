package maps

import (
	"github.com/google/uuid"
	"github.com/linkedin/goavro"
	"kafka-golang/lib/observable"
)

type schemaMapType map[uuid.UUID]*goavro.Codec

type SchemaMap struct {
	observable.LockableObservable
	Map schemaMapType
}

func (m SchemaMap) Insert(schemaUUID uuid.UUID, codec *goavro.Codec) bool {
	_, overwritten := m.Map[schemaUUID]
	m.Map[schemaUUID] = codec
	m.Notify(schemaUUID)
	return overwritten
}

func NewSchemaMap() SchemaMap {
	return SchemaMap{
		LockableObservable: observable.NewLockableObservable(),
		Map:                make(schemaMapType),
	}
}
