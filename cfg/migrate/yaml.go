package migrate

import (
	"bytes"
	"fmt"

	"gopkg.in/yaml.v2"
)

func Marshal(comment string, v any) ([]byte, error) {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "#%s\n", comment)
	encoder := yaml.NewEncoder(&buf)
	err := encoder.Encode(v)
	if err != nil {
		return nil, err
	}
	err = encoder.Close()
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
