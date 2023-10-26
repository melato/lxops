package cfg

import (
	"fmt"
	"os"

	"melato.org/lxops/yaml"
)

var Trace bool

const Comment = "#lxops-v1"

// MigrateFunc migrates a config file (represented by []byte) to another format.
type MigrateFunc func([]byte) ([]byte, error)

var ConfigFormats = make(map[string]MigrateFunc)

func SetMigrateFunc(comment string, fn MigrateFunc) {
	ConfigFormats[comment] = fn
}

/** Read raw config from yaml */
func ReadConfigYaml(file string) (*Config, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	if Trace {
		fmt.Printf("target config type: %s\n", Comment)
	}
	migrated := make(map[string]bool)
	for {
		comment := yaml.FirstLineComment(data)
		if Trace {
			fmt.Printf("config type: %s\n", comment)
		}
		if comment == Comment {
			break
		}
		if migrated[comment] {
			return nil, fmt.Errorf("config file migration loop: %s", comment)
		}
		migrated[comment] = true

		fn, supported := ConfigFormats[comment]
		if !supported {
			return nil, fmt.Errorf("no migration from config type %s", comment)
		}
		data, err = fn(data)
		if err != nil {
			return nil, fmt.Errorf("config type %s migration: %w", comment, err)
		}
	}
	var c Config
	err = yaml.Unmarshal(data, &c)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func PrintConfigYaml(config *Config) error {
	fmt.Printf("%s\n", Comment)
	data, err := yaml.Marshal(config)
	if err != nil {
		return err
	}
	fmt.Println(string(data))
	return nil
}
