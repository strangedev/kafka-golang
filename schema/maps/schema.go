package maps

import (
	"github.com/google/uuid"
	"github.com/linkedin/goavro"
	"github.com/strangedev/kafka-golang/lib"
)

type schemaMapType map[uuid.UUID]*goavro.Codec

// SchemaMap is a KeyObservable map of UUIDs to Avro codecs
type SchemaMap struct {
	lib.ConcurrentObservable
	Map schemaMapType
}

// Upsert inserts or updates a UUID, Codec pair into the map.
// All of the map's observers are notified of this change.
// It returns true, if the map entry did already exist and was overwritten.
func (m SchemaMap) Upsert(schemaUUID uuid.UUID, codec *goavro.Codec) bool {
	m.DataLock.Lock()
	_, overwritten := m.Map[schemaUUID]
	m.Map[schemaUUID] = codec
	m.DataLock.Unlock()
	m.Notify(schemaUUID)
	return overwritten
}

// NewSchemaMap constructs an empty SchemaMap with no observers.
func NewSchemaMap() SchemaMap {
	return SchemaMap{
		ConcurrentObservable: lib.NewConcurrentObservable(),
		Map:                  make(schemaMapType),
	}
}
