package lxops

import (
	"fmt"
	"os"

	"melato.org/lxops/yaml"
)

const Comment = "#lxdops"

/** Read raw config from yaml */
func ReadConfigYaml(file string) (*Config, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	if !yaml.FirstLineIs(data, Comment) {
		return nil, fmt.Errorf("%s: first line should be: %s\n", file, Comment)
	}
	var x Config
	err = yaml.Unmarshal(data, &x)
	if err != nil {
		return nil, err
	}
	return &x, err
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
