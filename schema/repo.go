package schema

import (
	"github.com/google/uuid"
)

type Repo interface {
	DecodeAvro(schema uuid.UUID, datum []byte) (interface{}, error)
	EncodeAvro(schema uuid.UUID, datum interface{}) ([]byte, error)
	WaitSchemaReady(schema uuid.UUID) chan bool
	WaitAliasReady(alias Alias) chan bool
	ListSchemata() []uuid.UUID
	ListAliases() []Alias
	WhoIs(alias Alias) (uuid.UUID, bool)
	GetSpecification(schema uuid.UUID) (specification string, ok bool)
	Count() int
}
