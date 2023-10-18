package property

import (
	"fmt"
	"strings"

	"melato.org/lxops/yaml"
)

type Options struct {
	PropertyKeyFiles []string `name:"f" usage:"yaml file"`
	KeyValues        []string `name:"D" usage:"key=value - set a property"`
}

func (t *Options) SetFile(properties Properties, keys []string, file string) error {
	var fileProperties Properties
	err := yaml.ReadFile(file, &fileProperties)
	if err != nil {
		return err
	}
	if len(keys) == 0 {
		for key, value := range fileProperties {
			properties[key] = value
		}
		return nil
	}
	return properties.Set(keys, fileProperties)
}

// ParseKeyValue - parse a string of the form <key1[.key2]...>=<value>
// The keys are separated by dots
// The first returned argument is the keys, and the second is the value
// If there is no "=", then the keys are nil and the value is the input string
func (t *Options) ParseKeyValue(keyValue string) ([]string, string) {
	kv := strings.SplitN(keyValue, "=", 2)
	if len(kv) != 2 {
		return nil, keyValue
	}
	compoundKey, value := kv[0], kv[1]
	keys := strings.Split(compoundKey, ".")
	return keys, value
}

func (t *Options) Apply(properties Properties) error {
	for _, kvFile := range t.PropertyKeyFiles {
		keys, file := t.ParseKeyValue(kvFile)
		err := t.SetFile(properties, keys, file)
		if err != nil {
			return err
		}
	}
	for _, kv := range t.KeyValues {
		keys, value := t.ParseKeyValue(kv)
		if len(keys) == 0 {
			return fmt.Errorf("missing key(s): %s", kv)
		}
		err := properties.Set(keys, value)
		if err != nil {
			return err
		}
	}
	return nil

}
