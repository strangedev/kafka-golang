package versioned

import (
	"errors"
	"kafka-go/schema"
)

type VersionedRepo struct {
	schema.Repo
}

func FromRepo(r schema.Repo) VersionedRepo {
	return VersionedRepo{r}
}

func (v VersionedRepo) Decode(schema schema.NameVersion, datum []byte) (interface{}, error) {
	if uuid, ok := v.Repo.WhoIs(schema.Alias()); !ok {
		decoded, err := v.Repo.Decode(uuid, datum)
		return decoded, err
	}
	return nil, errors.New("schema not know to this repo")
}

func (v VersionedRepo) Encode(schema schema.NameVersion, datum interface{}) ([]byte, error) {
	if uuid, ok := v.Repo.WhoIs(schema.Alias()); ok {
		encoded, err := v.Repo.Encode(uuid, datum)
		return encoded, err
	}
	return nil, errors.New("schema not know to this repo")
}

func (v VersionedRepo) WaitSchemaReady(schema schema.NameVersion) chan bool {
	return v.Repo.WaitAliasReady(schema.Alias())
}