// Package explorer contains things needed for the RESTful schema repository aggregate called the "explorer".
package explorer

import (
	"github.com/google/uuid"
	"github.com/strangedev/kafka-golang/schema"
)

// SchemaDTO is used by the explorer to encode its response body.
type SchemaDTO struct {
	UUID          uuid.UUID `json:"uuid"`
	Specification string    `json:"spec"`
}

// SchemataDTO is used by the explorer to encode its response body.
type SchemataDTO struct {
	Schemata []SchemaDTO `json:"schemata"`
}

// SchemaListDTO is used by the explorer to encode its response body.
type SchemaListDTO struct {
	Count    int         `json:"count"`
	Schemata []uuid.UUID `json:"schemata"`
}

// AliasDTO is used by the explorer to encode its response body.
type AliasDTO struct {
	Alias schema.Alias `json:"alias"`
	UUID  uuid.UUID    `json:"uuid"`
}

// AliasListDTO is used by the explorer to encode its response body.
type AliasListDTO struct {
	Aliases []schema.Alias `json:"aliases"`
	Count   int            `json:"count"`
}

// AliasesDTO is used by the explorer to encode its response body.
type AliasesDTO struct {
	Aliases []AliasDTO `json:"aliases"`
}
