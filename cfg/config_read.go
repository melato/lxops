package cfg

import (
	"fmt"

	"melato.org/configfile"
)

type ConfigReader struct {
	getVariable GetVariable

	Warn    bool
	Verbose bool
	file    string
	warned  bool
}

func NewConfigReader(getVariable GetVariable) *ConfigReader {
	return &ConfigReader{getVariable: getVariable}
}

func (r *ConfigReader) Read(file string) (*Config, error) {
	r.warned = false
	r.file = file
	if r.Verbose {
		r.warned = true
	}
	c, err := configfile.LoadFile[Config](file, func(data []byte, v any) error {
		switch c := v.(type) {
		case *Config:
			return Unmarshal(data, c)
		default:
			return fmt.Errorf("unexpected type: %T", v)
		}
	})
	if err != nil {
		return nil, err
	}
	return c, nil
}

// ReadConfig read a config file, and its included files
func ReadConfig(file string, getVariable GetVariable) (*Config, error) {
	r := NewConfigReader(getVariable)
	c, err := r.Read(file)
	if err != nil {
		return nil, err
	}
	c.File = file
	return c, nil
}
