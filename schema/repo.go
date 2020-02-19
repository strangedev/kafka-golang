package schema

import (
	"github.com/google/uuid"
)

type Repo interface {
	DecodeAvro(schema uuid.UUID, datum []byte) (interface{}, error)
	EncodeAvro(schema uuid.UUID, datum interface{}) ([]byte, error)
	WaitSchemaReady(schema uuid.UUID) chan bool
	ListSchemata() []uuid.UUID
	GetSpecification(schema uuid.UUID) (specification string, ok bool)
	Count() int
}

type AliasRepo interface {
	WhoIs(alias Alias) (uuid.UUID, bool)
	WaitAliasReady(alias Alias) chan bool
	ListAliases() []Alias
}

type VersionedRepo interface {
	Decode(schema NameVersion, datum []byte) (interface{}, error)
	Encode(schema NameVersion, datum interface{}) ([]byte, error)
	WaitSchemaVersionReady(schema NameVersion) chan bool
}
