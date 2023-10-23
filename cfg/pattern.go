package cfg

// Pattern is a string that is converted via property substitution, before it is used.
type Pattern string

type PatternSubstitution interface {
	Substitute(string) (string, error)
}

func (pattern Pattern) Substitute(properties PatternSubstitution) (string, error) {
	return properties.Substitute(string(pattern))
}

func (pattern Pattern) IsEmpty() bool {
	return string(pattern) == ""
}
