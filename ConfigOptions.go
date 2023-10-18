package lxops

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"
)

type ConfigContext interface {
	CurrentProject() string
}

type ConfigOptions struct {
	Project    string   `name:"project" usage:"the LXD project to use.  Overrides Config.Project"`
	Name       string   `name:"name" usage:"The name of the container to launch or configure.  If missing, use a separate container for each config, using the name of the config."`
	Properties []string `name:"P" usage:"a command-line property in the form <key>=<value>.  Command-line properties override instance and global properties"`
	properties map[string]string
	PropertyOptions
}

func (t *ConfigOptions) Init() error {
	return t.PropertyOptions.Init()
}

func (t *ConfigOptions) ConfigureProject(client ConfigContext) {
	if t.Project == "" {
		t.Project = client.CurrentProject()
	}
}

func (t *ConfigOptions) Configured() error {
	return t.PropertyOptions.Configured()
}

func (t *ConfigOptions) initProperties() error {
	if t.properties != nil {
		return nil
	}
	t.properties = make(map[string]string)
	for _, property := range t.Properties {
		i := strings.Index(property, "=")
		if i < 0 {
			return errors.New("missing value from property: " + property)
		}
		t.properties[property[0:i]] = property[i+1:]
	}
	return nil
}

func (t *ConfigOptions) UpdateConfig(config *Config) {
	if t.Project != "" {
		config.Project = t.Project
	}
	for key, value := range t.properties {
		if config.Properties == nil {
			config.Properties = make(map[string]string)
		}
		config.Properties[key] = value
	}
}

func (t *ConfigOptions) ReadConfig(file string) (*Config, error) {
	err := t.initProperties()
	if err != nil {
		return nil, err
	}
	var config *Config
	config, err = ReadConfig(file)
	if err != nil {
		return nil, err
	}

	if !config.Verify() {
		return nil, errors.New("invalid config")
	}
	t.UpdateConfig(config)
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
	return newInstance(t.GlobalProperties, config, name, includeSource)
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
