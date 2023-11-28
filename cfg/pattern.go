package cfg

import (
	"fmt"
)

// Pattern is a string that is converted via property substitution, before it is used.
type Pattern string

type PatternSubstitution interface {
	Substitute(string) (string, error)
}

func (pattern Pattern) Substitute(properties PatternSubstitution) (string, error) {
	result, err := properties.Substitute(string(pattern))
	if err != nil {
		return "", err
	}
	if Trace {
		fmt.Printf("substitute %s -> %s\n", pattern, result)
	}
	return result, nil
}

func (pattern Pattern) IsEmpty() bool {
	return string(pattern) == ""
}
