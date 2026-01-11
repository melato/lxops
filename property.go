package lxops

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"melato.org/lxops/cfg"
	"melato.org/lxops/util"
	"melato.org/lxops/yaml"
)

type PropertyOptions struct {
	PropertiesFile   string   `name:"properties" usage:"a file containing global config properties"`
	Properties       []string `name:"P" usage:"a command-line property in the form <key>=<value>.  Command-line properties override instance and global properties"`
	cliProperties    map[string]string
	GlobalProperties map[string]string `name:"-"`
	userProperties   map[string]string
}

func (t *PropertyOptions) Init() error {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return err
	}
	t.PropertiesFile = filepath.Join(configDir, "lxops", "properties.yaml")
	return nil
}

func (t *PropertyOptions) GetProperty(name string) (string, bool) {
	value, found := t.cliProperties[name]
	if !found {
		value, found = t.GlobalProperties[name]
	}
	if value == "" {
		found = false
	}
	return value, found
}

func ReadPropertiesDir(dir string, properties map[string]string) error {
	_, err := os.Stat(dir)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil
		}
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	files := make([]string, 0, len(entries))
	for _, entry := range entries {
		if !entry.IsDir() {
			name := entry.Name()
			if strings.HasSuffix(name, ".yaml") {
				files = append(files, name)
			}
		}
	}
	sort.Strings(files)
	for _, file := range files {
		file = filepath.Join(dir, file)
		err = yaml.ReadFile(file, properties)
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *PropertyOptions) loadPropertiesDir() error {
	t.GlobalProperties = make(map[string]string)
	configDir, err := os.UserConfigDir()
	if err == nil {
		propertiesDir := filepath.Join(configDir, "lxops", "properties.d")
		err = ReadPropertiesDir(propertiesDir, t.GlobalProperties)
	}
	return err
}

func (t *PropertyOptions) loadPropertiesFile() error {
	t.userProperties = make(map[string]string)
	if t.PropertiesFile != "" {
		_, err := os.Stat(t.PropertiesFile)
		if err != nil {
			return err
		}
		err = yaml.ReadFile(t.PropertiesFile, t.userProperties)
		if err != nil {
			return err
		}
		for name, value := range t.userProperties {
			t.GlobalProperties[name] = value
		}
	}
	return nil
}

func (t *PropertyOptions) parseCliProperties() error {
	t.cliProperties = make(map[string]string)
	for _, property := range t.Properties {
		i := strings.Index(property, "=")
		if i < 0 {
			return fmt.Errorf("missing value from property: %s", property)
		}
		t.cliProperties[property[0:i]] = property[i+1:]
	}
	return nil
}

func (t *PropertyOptions) Configured() error {
	err := t.loadPropertiesDir()
	if err == nil {
		err = t.loadPropertiesFile()
	}
	if err == nil {
		err = t.parseCliProperties()
	}
	return err
}

func (t *PropertyOptions) List() {
	util.PrintMap(t.GlobalProperties)
}

func (t *PropertyOptions) File() {
	fmt.Println(t.PropertiesFile)
}

func (t *PropertyOptions) Set(key, value string) error {
	if t.userProperties == nil {
		t.userProperties = make(map[string]string)
	}
	t.userProperties[key] = value
	t.GlobalProperties[key] = value
	if t.PropertiesFile != "" {
		dir := filepath.Dir(t.PropertiesFile)
		err := os.MkdirAll(dir, os.FileMode(0775))
		if err != nil {
			return err
		}
		return yaml.WriteFile(t.userProperties, t.PropertiesFile)
	}
	return nil
}

func (t *PropertyOptions) Get(key string) error {
	value := t.GlobalProperties[key]
	fmt.Printf("%s\n", value)
	return nil
}

func (t *PropertyOptions) NewConfigReader() *cfg.ConfigReader {
	return cfg.NewConfigReader(t.GetProperty)
}
