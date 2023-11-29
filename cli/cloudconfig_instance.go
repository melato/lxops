package cli

import (
	"fmt"
)

// CloudconfigInstanceOps - runs cloudconfig files on one instance
type CloudconfigInstanceOps struct {
	Cloudconfig *Cloudconfig
	Instance    string `name:"i" usage:"instance to configure"`
}

func (t *CloudconfigInstanceOps) Configured() error {
	if t.Instance == "" {
		return fmt.Errorf("missing instance")
	}
	return t.Cloudconfig.Configured()
}

func (t *CloudconfigInstanceOps) Apply(configFiles ...string) error {
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
