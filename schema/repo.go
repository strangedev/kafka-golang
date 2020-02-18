package schema

import (
	"github.com/google/uuid"
)

type Repo interface {
	Decode(schema uuid.UUID, datum []byte) (interface{}, error)
	Encode(schema uuid.UUID, datum interface{}) ([]byte, error)
	WaitSchemaReady(schema uuid.UUID) chan bool
	WaitAliasReady(alias Alias) chan bool
	ListSchemata() []uuid.UUID
	ListAliases() []Alias
	WhoIs(alias Alias) (uuid.UUID, bool)
	GetSpecification(schema uuid.UUID) (specification string, ok bool)
	Count() int
	Run() (stop chan bool)
}

type VersionedRepo interface {
	Decode(schema NameVersion, datum []byte) (interface{}, error)
	Encode(schema NameVersion, datum interface{}) ([]byte, error)
	WaitSchemaReady(schema NameVersion) chan bool
}
