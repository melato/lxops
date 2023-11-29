package cli

import (
	"fmt"
	"io"
	"os"

	"melato.org/cloudconfig"
)

// CloudconfigFileOps - runs one cloudconfig file to multiple instances
type CloudconfigFileOps struct {
	Cloudconfig *Cloudconfig
	File        string `name:"f" usage:"cloudconfig file"`
	cc          *cloudconfig.Config
}

func (t *CloudconfigFileOps) Configured() error {
	err := t.Cloudconfig.Configured()
	if err != nil {
		return err
	}
	var data []byte
	if t.File == "" {
		data, err = io.ReadAll(os.Stdin)
	} else {
		data, err = os.ReadFile(t.File)
	}
	if err != nil {
		return err
	}
	t.cc, err = cloudconfig.Unmarshal(data)
	if err != nil {
		return err
	}
	return nil
}

func (t *CloudconfigFileOps) ApplyInstance(instance string) error {
	configurer, err := t.Cloudconfig.NewConfigurer(instance)
	if err != nil {
		return err
	}
	return configurer.Apply(t.cc)
}

func (t *CloudconfigFileOps) Apply(instances ...string) error {
	for _, instance := range instances {
		fmt.Printf("%s:\n", instance)
		err := t.ApplyInstance(instance)
		if err != nil {
			return err
		}
	}
	return nil
}
