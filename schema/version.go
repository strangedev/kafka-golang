package schema

import "fmt"

type NameVersion struct {
	Name    string
	Version uint
}

type Version interface {
	IsOrigin() bool
	GetVersion() uint
	GetName() string
	GetPrevious() Version
	GetNext() Version
	String() string
	Alias() Alias
}

func NewVersionOrigin(alias string) NameVersion {
	return NameVersion{Name: alias, Version: 0}
}

func (v NameVersion) IsOrigin() bool {
	return v.Version == 0
}

func (v NameVersion) GetVersion() uint {
	return v.Version
}

func (v NameVersion) GetName() string {
	return v.Name
}

func (v NameVersion) GetPrevious() Version {
	if v.Version == 0 {
		return v
	}
	previousVersion := v.Version - 1
	return NameVersion{Name: v.Name, Version: previousVersion}
}

func (v NameVersion) GetNext() Version {
	previousVersion := v.Version + 1
	return NameVersion{Name: v.Name, Version: previousVersion}
}

func (v NameVersion) String() string {
	return fmt.Sprintf("%s-v%x", v.Name, v.Version)
}

func (v NameVersion) Alias() Alias {
	return Alias(v.String())
}
