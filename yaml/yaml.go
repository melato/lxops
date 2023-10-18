// convenience utility wrappers for yaml
package yaml

import (
	"fmt"
	"io"
	"os"

	"gopkg.in/yaml.v2"
)

/** read a yaml object from a file
    Return the object for convenience, so that you can return both the object and the error,

	example:
	var conf Conf
	return script.ReadYamlFile(&conf, file) */

func ReadFile(file string, v interface{}) error {
	data, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(data, v)
	if err != nil {
		return fmt.Errorf("%s: %v", file, err)
	}
	return nil
}

func Marshal(v interface{}) ([]byte, error) {
	return yaml.Marshal(v)
}
func Unmarshal(data []byte, v interface{}) error {
	return yaml.Unmarshal(data, v)
}

func WriteFile(v interface{}, file string) error {
	data, err := yaml.Marshal(v)
	if err != nil {
		return err
	}
	return os.WriteFile(file, data, 0664)
}

func Write(writer io.Writer, v interface{}) error {
	encoder := yaml.NewEncoder(writer)
	return encoder.Encode(v)
}

func Print(v interface{}) error {
	return Write(os.Stdout, v)
}
