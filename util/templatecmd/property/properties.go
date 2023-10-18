package property

import (
	"fmt"
	"reflect"
	"strings"
)

// Properties -- a map that holds generic properties that are used in templates.
// properties are hierarchical, so that a property value can be another map[any]any
// The keys are typically strings, but Properties is defined as a map from any Go type to any Go type,
// so that it can be loaded generically from a yaml file.
type Properties map[any]any

func (t Properties) Set(keys []string, value any) error {
	p := t
	for i, k := range keys {
		if i == len(keys)-1 {
			p[k] = value
		} else {
			v, found := p[k]
			var kMap map[any]any
			if found {
				var isMap bool
				kMap, isMap = v.(map[any]any)
				if !isMap {
					kProperties, isProperties := v.(Properties)
					if isProperties {
						kMap = map[any]any(kProperties)
					} else {
						key1 := strings.Join(keys[0:i+1], ".")
						return fmt.Errorf("not a map: key=%s type=%v", key1, reflect.TypeOf(v))
					}
				}
			} else {
				kMap = make(map[any]any)
			}
			p[k] = kMap
			p = kMap
		}
	}
	return nil
}
