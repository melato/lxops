package util

import (
	"fmt"
	"regexp"
)

var reParen = regexp.MustCompile(`\([^()]*\)`)

type substitution struct {
	Error  error
	Lookup func(string) (string, bool)
}

func (t *substitution) Get(key string) string {
	value, found := t.Lookup(key)
	if found {
		return value
	}
	if t.Error == nil {
		t.Error = fmt.Errorf("no such key: %s", key)
	}
	return ""
}

func (t *substitution) Replace(key string) string {
	key = key[1 : len(key)-1]
	return t.Get(key)
}

func Substitute(pattern string, lookup func(string) (string, bool)) (string, error) {
	if pattern == "" {
		return "", nil
	}
	var sub substitution
	sub.Lookup = lookup
	value := reParen.ReplaceAllStringFunc(pattern, sub.Replace)
	if sub.Error != nil {
		return "", sub.Error
	}
	return value, nil
}
