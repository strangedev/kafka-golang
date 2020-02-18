package schema

type Alias string

func (a Alias) String() string {
	return string(a)
}
