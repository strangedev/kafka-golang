// Package command contains commands for updating the schema stored in Kafka.
// Since the schema repository can't use itself to specify the Command/Message format
// for update actions, it uses plain JSON.
package command

import "github.com/google/uuid"

// UpdateRequest sets the given UUID to equal the given plain-text Avro spec.
type UpdateRequest struct {
	UUID uuid.UUID `json:"UUID"`
	Spec string    `json:"spec"`
}

// AliasRequest sets the given Alias to equal the given UUID.
type AliasRequest struct {
	UUID  uuid.UUID `json:"UUID"`
	Alias string    `json:"alias"`
}
