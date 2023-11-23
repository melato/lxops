package util

import (
	"fmt"

	"melato.org/lxops/template"
)

/*
* Uses template.Paren and a list of map[string]string to implement pattern substitution on strings,
Replaces parenthesized expressions as follows:
(key) ->
*/
type CascadingProperties struct {
	Maps []map[string]string
}

// AddMap adds a map to the list of maps used to lookup keys.
// maps are checked in the order in which they were added,
// so the given map is used only for keys that are not in the previous maps.
func (t *CascadingProperties) AddMap(m map[string]string) {
	t.Maps = append(t.Maps, m)
}

func (t *CascadingProperties) Get(key string) (string, error) {
	for _, properties := range t.Maps {
		value, found := properties[key]
		if found {
			return value, nil
		}
	}
	return "", fmt.Errorf("no such key: %s", key)
}

func (t *CascadingProperties) Substitute(pattern string) (string, error) {
	if pattern == "" {
		return "", nil
	}
	tpl, err := template.Paren.NewTemplate(pattern)
	if err != nil {
		return "", err
	}
	return tpl.Applyf(t.Get)
}
