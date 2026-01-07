package util

/*
* Uses Paren and a list of map[string]string to implement pattern substitution on strings,
Replaces parenthesized expressions as follows:
(key) -> with the first value it finds in one of the maps, searching them in order.
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

func (t *CascadingProperties) Lookup(key string) (string, bool) {
	for _, properties := range t.Maps {
		value, found := properties[key]
		if found {
			return value, true
		}
	}
	return "", false
}

func (t *CascadingProperties) Substitute(pattern string) (string, error) {
	return Substitute(pattern, t.Lookup)
}
