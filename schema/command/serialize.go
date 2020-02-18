package command

import "github.com/google/uuid"

type UpdateRequest struct {
	UUID uuid.UUID `json:"UUID"`
	Spec string    `json:"spec"`
}

type AliasRequest struct {
	UUID  uuid.UUID `json:"UUID"`
	Alias string    `json:"alias"`
}