package lxops

import (
	"os"

	"melato.org/cloudconfig"
	"melato.org/lxops/cfg"
	"melato.org/lxops/srv"
)

type Configurer struct {
	Client srv.Client `name:"-"`
	ConfigOptions
	Trace bool `name:"trace,t" usage:"print exec arguments"`
	DryRunFlag
}

func (t *Configurer) Init() error {
	return t.ConfigOptions.Init()
}

func (t *Configurer) Configured() error {
	t.ConfigOptions.ConfigureProject(t.Client)
	return t.ConfigOptions.Configured()
}

/** run things inside the container:  install packages, create users, run scripts */
func (t *Configurer) ConfigureContainer(instance *Instance) error {
	config := instance.Config
	container := instance.Container()
	server, err := t.Client.ProjectInstanceServer(config.Project)
	if err != nil {
		return err
	}
	if !t.DryRun {
		err := server.WaitForNetwork(container)
		if err != nil {
			return err
		}
	}
	if len(config.CloudConfigFiles) > 0 {
		base, err := server.NewConfigurer(instance.Name)
		base.SetLogWriter(os.Stdout)
		if err != nil {
			return err
		}
		configurer := cloudconfig.NewConfigurer(base)
		configurer.OS, err = cfg.OSType(config.Ostype)
		if err != nil {
			return err
		}
		configurer.Log = os.Stdout
		files := make([]string, len(config.CloudConfigFiles))
		for i, file := range config.CloudConfigFiles {
			files[i] = string(file)
		}
		err = configurer.ApplyConfigFiles(files...)
		if err != nil {
			return err
		}
	}
	return nil
}
