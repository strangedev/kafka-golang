package explorer

import (
	"github.com/google/uuid"
	"github.com/strangedev/kafka-golang/schema"
)

type SchemaDTO struct {
	UUID          uuid.UUID `json:"uuid"`
	Specification string    `json:"spec"`
}

type SchemataDTO struct {
	Schemata []SchemaDTO `json:"schemata"`
}

type SchemaListDTO struct {
	Count    int         `json:"count"`
	Schemata []uuid.UUID `json:"schemata"`
}

type AliasDTO struct {
	Alias schema.Alias `json:"alias"`
	UUID  uuid.UUID    `json:"uuid"`
}

type AliasListDTO struct {
	Aliases []schema.Alias `json:"aliases"`
	Count   int            `json:"count"`
}

type AliasesDTO struct {
	Aliases []AliasDTO `json:"aliases"`
}
