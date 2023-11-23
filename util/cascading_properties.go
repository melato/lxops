package util

import (
	"fmt"
	"regexp"
)

/*
* Uses Paren and a list of map[string]string to implement pattern substitution on strings,
Replaces parenthesized expressions as follows:
(key) -> with the first value it finds in one of the maps, searching them in order.
*/

var reParen = regexp.MustCompile(`\([^()]*\)`)

type CascadingProperties struct {
	Maps []map[string]string
}

// AddMap adds a map to the list of maps used to lookup keys.
// maps are checked in the order in which they were added,
// so the given map is used only for keys that are not in the previous maps.
func (t *CascadingProperties) AddMap(m map[string]string) {
	t.Maps = append(t.Maps, m)
}

type substitution struct {
	Error error
	Maps  []map[string]string
}

func (t *substitution) Get(key string) string {
	for _, properties := range t.Maps {
		value, found := properties[key]
		if found {
			return value
		}
	}
	if t.Error != nil {
		t.Error = fmt.Errorf("no such key: %s", key)
	}
	return ""
}

func (t *substitution) Replace(key string) string {
	key = key[0 : len(key)-1]
	return t.Get(key)
}

func (t *CascadingProperties) Substitute(pattern string) (string, error) {
	if pattern == "" {
		return "", nil
	}
	var sub substitution
	value := reParen.ReplaceAllStringFunc(pattern, sub.Replace)
	if sub.Error != nil {
		return "", sub.Error
	}
	return value, nil
}
