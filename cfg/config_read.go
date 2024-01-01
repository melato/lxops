package cfg

import (
	"fmt"
	"os"
	"path/filepath"
)

type ConfigReader struct {
	Warn     bool
	Verbose  bool
	included map[string]bool
	file     string
	warned   bool
}

func (r *ConfigReader) isIncluded(file string) bool {
	return r.included[file]
}

func (r *ConfigReader) addIncluded(file string) {
	if r.included == nil {
		r.included = make(map[string]bool)
	}
	r.included[file] = true
}

func (t *OS) Merge(c *OS) error {
	if c == nil {
		// keep the one we have
		return nil
	}
	if t.Name != c.Name {
		return fmt.Errorf("cannot merge incompatible OSs: %s, %s", t.Name, c.Name)
	}
	return nil
}

func (r *ConfigReader) warn(format string, arg ...interface{}) {
	if !r.warned {
		fmt.Println(r.file)
		r.warned = true
	}
	fmt.Printf(format, arg...)
}

// mergeMaps add map entries from b to a.  b overrides a
func (r *ConfigReader) mergeMaps(a, b map[string]string) (map[string]string, error) {
	if a == nil && b == nil {
		return nil, nil
	}
	if a == nil {
		a = make(map[string]string)
	}
	for key, value := range b {
		if r.Warn {
			oldValue, _ := a[key]
			if oldValue != value && oldValue != "" {
				fmt.Fprintf(os.Stderr, "%s: \"%s\" overrides \"%s\"\n", key, value, oldValue)
			}
		}
		a[key] = value
	}
	return a, nil
}

func (r *ConfigReader) mergeSource(t, c *Source) {
	if c.Origin != "" {
		t.Origin = c.Origin
	}
	if c.DeviceTemplate != "" {
		t.DeviceTemplate = c.DeviceTemplate
	}
	if c.DeviceOrigin != "" {
		t.DeviceOrigin = c.DeviceOrigin
	}
	if c.SourceConfig != "" {
		t.SourceConfig = c.SourceConfig
	}
}

func (r *ConfigReader) mergeInherit(t, c *ConfigInherit) error {
	if c.Project != "" {
		t.Project = c.Project
	}
	if c.Profile != "" {
		t.Profile = c.Profile
	}
	if c.DeviceOwner != "" {
		t.DeviceOwner = c.DeviceOwner
	}
	var err error
	t.Properties, err = r.mergeMaps(t.Properties, c.Properties)
	if err != nil {
		return err
	}
	if err != nil {
		return err
	}

	r.mergeSource(&t.Source, &c.Source)

	if len(c.LxcOptions) != 0 {
		t.LxcOptions = c.LxcOptions
	}

	if t.Filesystems == nil {
		t.Filesystems = make(map[string]*Filesystem)
	}
	for id, fs := range c.Filesystems {
		if r.Warn {
			_, exists := t.Filesystems[id]
			if exists {
				fmt.Printf("filesystem %s is overriden\n", id)
			}
		}
		t.Filesystems[id] = fs
	}
	if t.Devices == nil {
		t.Devices = make(map[string]*Device)
	}
	for id, d := range c.Devices {
		if r.Warn {
			_, exists := t.Devices[id]
			if exists {
				fmt.Printf("device %s is overriden\n", id)
			}
		}
		t.Devices[id] = d
	}

	t.Profiles = append(t.Profiles, c.Profiles...)
	t.CloudConfigFiles = append(t.CloudConfigFiles, c.CloudConfigFiles...)
	return nil
}

func (r *ConfigReader) mergeFile(t *Config, file string) error {
	if r.isIncluded(file) {
		fmt.Fprintf(os.Stderr, "ignoring duplicate include: %s\n", file)
		return nil
	}
	if r.Verbose {
		fmt.Println(file)
	}
	config, err := ReadRawConfig(file)
	if err != nil {
		return err
	}
	dir := filepath.Dir(file)
	config.ResolvePaths(dir)
	if len(r.included) == 0 {
		t.ConfigTop = config.ConfigTop
	}
	r.addIncluded(file)
	if config.Ostype != "" {
		t.Ostype = config.Ostype
	}
	if config.Image != "" {
		t.Image = config.Image
	}
	for _, f := range config.Include {
		err := r.mergeFile(t, string(f))
		if err != nil {
			return err
		}
	}
	return r.mergeInherit(&t.ConfigInherit, &config.ConfigInherit)
}

func (r *ConfigReader) Read(file string) (*Config, error) {
	r.warned = false
	r.included = nil
	r.file = file
	if r.Verbose {
		r.warned = true
	}
	result := &Config{}
	err := r.mergeFile(result, file)
	if err != nil {
		return nil, err
	}
	return result, err
}

// ReadConfig read a config file, and its included files
func ReadConfig(file string) (*Config, error) {
	r := &ConfigReader{}
	c, err := r.Read(file)
	if err != nil {
		return nil, err
	}
	c.File = file
	return c, nil
}
