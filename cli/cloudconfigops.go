package cli

import (
	"fmt"
	"io"
	"os"

	"melato.org/cloudconfig"
)

// CloudconfigOps - runs one cloudconfig file to multiple instances
type CloudconfigOps struct {
	Cloudconfig *Cloudconfig
	Instance    string `name:"i" usage:"instance"`
	File        string `name:"f" usage:"cloudconfig file"`
}

func (t *CloudconfigOps) Configured() error {
	return t.Cloudconfig.Configured()
}

func (t *CloudconfigOps) readFile() (*cloudconfig.Config, error) {
	var data []byte
	var err error
	if t.File == "" {
		data, err = io.ReadAll(os.Stdin)
	} else {
		data, err = os.ReadFile(t.File)
	}
	if err != nil {
		return nil, err
	}
	return cloudconfig.Unmarshal(data)
}

func (t *CloudconfigOps) ApplyFiles(configFiles ...string) error {
	configurer, err := t.Cloudconfig.NewConfigurer(t.Instance)
	if err != nil {
		return err
	}
	if len(configFiles) == 0 {
		return configurer.ApplyStdin()
	} else {
		return configurer.ApplyConfigFiles(configFiles...)
	}
}

func (t *CloudconfigOps) ApplyInstances(instances ...string) error {
	cc, err := t.readFile()
	if err != nil {
		return nil
	}
	for _, instance := range instances {
		fmt.Printf("%s:\n", instance)
		configurer, err := t.Cloudconfig.NewConfigurer(instance)
		if err != nil {
			return err
		}
		err = configurer.Apply(cc)
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *CloudconfigOps) Apply(args ...string) error {
	if t.Instance != "" && t.File == "" {
		return t.ApplyFiles(args...)
	}
	if t.File != "" && t.Instance == "" {
		return t.ApplyInstances(args...)
	}
	if t.File != "" && t.Instance != "" {
		if len(args) > 0 {
			return fmt.Errorf("cannot use -i, -f and arguments together")
		}
		return t.ApplyInstances(t.Instance)
	}
	return fmt.Errorf("at least -i or -f is required")
}
