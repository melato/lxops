package cli

import (
	"fmt"
	"os"

	"melato.org/cloudconfig"
	"melato.org/lxops/cfg"
	"melato.org/lxops/srv"
)

// Cloudconfig - cloudconfig operations
type Cloudconfig struct {
	Client srv.Client `name:"-"`
	OSType string     `name:"ostype" usage:"OS type.  Use ? for a list."`
	ostype cloudconfig.OSType
	server srv.InstanceServer `name:"-"`
}

func (t *Cloudconfig) Configured() error {
	if t.OSType == "?" {
		ListOSTypes()
		return fmt.Errorf("")
	}
	var err error
	if t.OSType != "" {
		t.ostype, err = cfg.OSType(t.OSType)
		if err != nil {
			return err
		}
	}

	server, err := t.Client.CurrentInstanceServer()
	if err != nil {
		return err
	}
	t.server = server
	return nil
}

func (t *Cloudconfig) NewConfigurer(instance string) (*cloudconfig.Configurer, error) {
	base, err := t.server.NewConfigurer(instance)
	if err != nil {
		return nil, err
	}
	base.SetLogWriter(os.Stdout)
	configurer := cloudconfig.NewConfigurer(base)
	configurer.OS = t.ostype
	configurer.Log = os.Stdout
	return configurer, nil
}
