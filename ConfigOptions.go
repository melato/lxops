package lxops

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"melato.org/lxops/cfg"
)

type ConfigContext interface {
	CurrentProject() string
}

type ConfigOptions struct {
	Project       string   `name:"project" usage:"the instance server project to use.  Overrides config project"`
	Name          string   `name:"name" usage:"The name of the instance.  If missing, use the base name of the config file"`
	Properties    []string `name:"P" usage:"a command-line property in the form <key>=<value>.  Command-line properties override instance and global properties"`
	cliProperties map[string]string
	PropertyOptions
}

func (t *ConfigOptions) Init() error {
	return t.PropertyOptions.Init()
}

func hasKey(m map[string]string, key string) bool {
	_, exists := m[key]
	return exists

}

func (t *ConfigOptions) HasProperty(name string) bool {
	return hasKey(t.cliProperties, name) || hasKey(t.GlobalProperties, name)
}

func (t *ConfigOptions) ConfigureProject(client ConfigContext) {
	if t.Project == "" {
		t.Project = client.CurrentProject()
	}
}

func (t *ConfigOptions) Configured() error {
	t.cliProperties = make(map[string]string)
	for _, property := range t.Properties {
		i := strings.Index(property, "=")
		if i < 0 {
			return fmt.Errorf("missing value from property: %s", property)
		}
		t.cliProperties[property[0:i]] = property[i+1:]
	}
	return t.PropertyOptions.Configured()
}

func (t *ConfigOptions) ReadConfig(file string) (*cfg.Config, error) {
	config, err := cfg.ReadConfig(file, t.HasProperty)
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
