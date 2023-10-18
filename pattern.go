package lxops

import (
	"melato.org/lxops/util"
)

// Pattern is a string that is converted via property substitution, before it is used.
type Pattern string

func (pattern Pattern) Substitute(properties *util.PatternProperties) (string, error) {
	return properties.Substitute(string(pattern))
}

func (pattern Pattern) IsEmpty() bool {
	return string(pattern) == ""
}
