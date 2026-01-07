package lxops

import (
	"errors"
	"fmt"
	"path/filepath"

	"melato.org/lxops/cfg"
)

type ConfigContext interface {
	CurrentProject() string
}

type ConfigOptions struct {
	PropertyOptions
	Project string `name:"project" usage:"the instance server project to use.  Overrides config project"`
	Name    string `name:"name" usage:"The name of the instance.  If missing, use the base name of the config file"`
}

func (t *ConfigOptions) Init() error {
	return t.PropertyOptions.Init()
}

func (t *ConfigOptions) ConfigureProject(client ConfigContext) {
	if t.Project == "" {
		t.Project = client.CurrentProject()
	}
}

func (t *ConfigOptions) ReadConfig(file string) (*cfg.Config, error) {
	config, err := cfg.ReadConfig(file, t.GetProperty)
	if err != nil {
		return nil, err
	}

	if !config.Verify() {
		return nil, errors.New("invalid config")
	}
	return config, nil
}

func BaseName(file string) string {
	name := filepath.Base(file)
	ext := filepath.Ext(name)
	return name[0 : len(name)-len(ext)]
}

func (t *ConfigOptions) instance2(file string, includeSource bool) (*Instance, error) {
	var name string
	if t.Name != "" {
		name = t.Name
	} else {
		name = BaseName(file)
	}
	config, err := t.ReadConfig(file)
	if err != nil {
		return nil, err
	}
	return newInstance(t.cliProperties, t.GlobalProperties, config, name, includeSource)
}

func (t *ConfigOptions) Instance(file string) (*Instance, error) {
	return t.instance2(file, false)
}

func (t *ConfigOptions) Instance2(file string, includeSource bool) (*Instance, error) {
	return t.instance2(file, includeSource)
}

func (t *ConfigOptions) RunInstances(f func(*Instance) error, includeSource bool, args ...string) error {
	if t.Name != "" && len(args) != 1 {
		return errors.New("--name can be used with only one config file")
	}
	for _, arg := range args {
		instance, err := t.Instance2(arg, includeSource)
		if err != nil {
			return err
		}
		err = f(instance)
		if err != nil {
			return fmt.Errorf("%s: %w", arg, err)
		}
	}
	return nil
}

func (t *ConfigOptions) InstanceFunc(f func(*Instance) error, includeSource bool) func(configs []string) error {
	return func(configs []string) error {
		return t.RunInstances(f, includeSource, configs...)
	}
}
