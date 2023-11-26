package cli

import (
	"fmt"
	"os"

	"melato.org/cloudconfig"
	"melato.org/lxops/cfg"
)

type CloudconfigOps struct {
	InstanceOps *InstanceOps `name:"-"`
	Instance    string       `name:"i" usage:"LXD instance to configure"`
	OSType      string       `name:"ostype" usage:"OS type"`
}

func (t *CloudconfigOps) Apply(configFiles ...string) error {
	var ostype cloudconfig.OSType
	var err error
	if t.OSType != "" {
		ostype, err = cfg.OSType(t.OSType)
		if err != nil {
			return err
		}
	}
	if t.Instance == "" {
		return fmt.Errorf("missing instance")
	}
	base, err := t.InstanceOps.server.NewConfigurer(t.Instance)
	if err != nil {
		return err
	}
	base.SetLogWriter(os.Stdout)
	configurer := cloudconfig.NewConfigurer(base)
	configurer.OS = ostype
	configurer.Log = os.Stdout
	if len(configFiles) == 0 {
		return configurer.ApplyStdin()
	} else {
		return configurer.ApplyConfigFiles(configFiles...)
	}
}
