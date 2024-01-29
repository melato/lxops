package lxops

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"melato.org/lxops/cfg"
	"melato.org/lxops/srv"
	"melato.org/lxops/util"
	"melato.org/script"
	"melato.org/table3"
)

type Instance struct {
	globalProperties map[string]string
	Config           *cfg.Config
	cliProperties    map[string]string
	Name             string
	profile          string
	containerSource  *ContainerSource
	deviceSource     *DeviceSource
	Properties       cfg.PatternSubstitution
	fspaths          map[string]*InstanceFS
	sourceConfig     *cfg.Config
}

func (t *Instance) substitute(e *error, pattern cfg.Pattern, defaultPattern cfg.Pattern) string {
	if pattern == "" {
		pattern = defaultPattern
	}
	value, err := pattern.Substitute(t.Properties)
	if err != nil {
		*e = err
	}
	return value
}

func (instance *Instance) createBuiltins() map[string]string {
	properties := make(map[string]string)
	name := instance.Name
	properties["instance"] = name
	project := instance.Config.Project
	var projectSlash, project_instance string
	if project == "" || project == "default" {
		project = "default"
		projectSlash = ""
		project_instance = name
	} else {
		projectSlash = project + "/"
		project_instance = project + "_" + name
	}
	properties["project"] = project
	properties["project/"] = projectSlash
	properties["project_instance"] = project_instance
	return properties
}

func mergeProperties(a, b map[string]string) {
	for k, v := range b {
		a[k] = v
	}
}

func (instance *Instance) EffectiveProperties() map[string]string {
	properties := make(map[string]string)
	mergeProperties(properties, instance.globalProperties)
	mergeProperties(properties, instance.Config.Properties)
	mergeProperties(properties, instance.cliProperties)
	return properties
}

func (instance *Instance) newProperties() cfg.PatternSubstitution {
	cascade := &util.CascadingProperties{}
	cascade.AddMap(instance.createBuiltins())
	cascade.AddMap(instance.cliProperties)
	cascade.AddMap(instance.Config.Properties)
	cascade.AddMap(instance.globalProperties)

	return cascade
}

func newInstance(cliProperties, globalProperties map[string]string, config *cfg.Config, name string, includeSource bool) (*Instance, error) {
	t := &Instance{
		globalProperties: globalProperties,
		cliProperties:    cliProperties,
		Config:           config,
		Name:             name}
	t.Properties = t.newProperties()
	var err error
	t.profile = t.substitute(&err, config.Profile, "(instance).lxops")
	if err != nil {
		return nil, err
	}
	if includeSource {
		t.containerSource, err = t.newContainerSource()
		if err != nil {
			return nil, err
		}
		t.deviceSource, err = t.newDeviceSource()
		if err != nil {
			return nil, err
		}
	} else {
		t.containerSource = &ContainerSource{}
		t.deviceSource = &DeviceSource{}
	}
	return t, nil
}

func NewInstance(cliProperties, globalProperties map[string]string, config *cfg.Config, name string) (*Instance, error) {
	return newInstance(cliProperties, globalProperties, config, name, true)
}

func (t *Instance) NewInstance(name string) (*Instance, error) {
	return NewInstance(nil, t.globalProperties, t.Config, name)
}

func (t *Instance) ContainerSource() *ContainerSource {
	return t.containerSource
}

func (t *Instance) DeviceSource() *DeviceSource {
	return t.deviceSource
}

func (t *Instance) ProfileName() string {
	return t.profile
}

// Container is the same as Name
func (t *Instance) Container() string {
	return t.Name
}

func (t *Instance) Filesystems() (map[string]*InstanceFS, error) {
	if t.fspaths == nil {
		fspaths := make(map[string]*InstanceFS)
		for id, fs := range t.Config.Filesystems {
			path, err := fs.Pattern.Substitute(t.Properties)
			if err != nil {
				return nil, err
			}
			fspaths[id] = &InstanceFS{Id: id, Path: path, Filesystem: fs}
		}
		t.fspaths = fspaths
	}
	return t.fspaths, nil
}

func (t *Instance) FilesystemList() ([]*InstanceFS, error) {
	paths, err := t.Filesystems()
	if err != nil {
		return nil, err
	}
	var list []*InstanceFS
	for _, path := range paths {
		list = append(list, path)
	}
	InstanceFSList(list).Sort()
	return list, nil
}

func (t *Instance) DeviceList() ([]InstanceDevice, error) {
	var devices []InstanceDevice
	for name, device := range t.Config.Devices {
		d := InstanceDevice{Name: name, Device: device}
		dir, err := t.DeviceDir(name, device)
		if err != nil {
			return nil, err
		}
		d.Source = dir
		devices = append(devices, d)
	}

	InstanceDeviceList(devices).Sort()
	return devices, nil
}

func (t *Instance) DeviceDir(deviceId string, device *cfg.Device) (string, error) {
	dir, err := device.Dir.Substitute(t.Properties)
	if err != nil {
		return "", err
	}
	if strings.HasPrefix(dir, "/") {
		return dir, nil
	}
	if device.Filesystem == "" {
		return "", fmt.Errorf("device %s has no filesystem but relative dir: %s", deviceId, dir)
	}
	if dir == "" {
		dir = deviceId
	} else if device.Dir == "." {
		dir = ""
	}

	fspaths, err := t.Filesystems()
	if err != nil {
		return "", err
	}
	fsPath, exists := fspaths[device.Filesystem]
	if !exists {
		return "", nil
	}

	if dir != "" {
		return filepath.Join(fsPath.Dir(), dir), nil
	} else {
		return fsPath.Dir(), nil
	}
}

// Snapshot creates a snapshot of all ZFS filesystems of the instance
func (instance *Instance) Snapshot(name string) error {
	filesystems, err := instance.FilesystemList()
	if err != nil {
		return err
	}
	s := &script.Script{Trace: true}
	for _, fs := range filesystems {
		if fs.IsZfs() && !fs.Filesystem.Transient {
			s.Run("sudo", "zfs", "snapshot", fs.Path+"@"+name)
		}
	}
	return s.Error()
}

// Rollback calls zfs rollback -r on the non-transient ZFS filesystems of the instance
func (instance *Instance) Rollback(name string) error {
	filesystems, err := instance.FilesystemList()
	if err != nil {
		return err
	}
	s := &script.Script{Trace: true}
	for _, fs := range filesystems {
		if fs.IsZfs() && !fs.Filesystem.Transient {
			s.Run("sudo", "zfs", "rollback", "-r", fs.Path+"@"+name)
		}
	}
	return s.Error()
}

// GetSourceConfig returns the parsed configuration specified by Config.SourceConfig
// If there is no Config.SourceConfig, it returns this instance's config
// It returns a non nil *Config or an error.
func (t *Instance) GetSourceConfig() (*cfg.Config, error) {
	if t.Config.SourceConfig == "" {
		return t.Config, nil
	}
	if t.sourceConfig == nil {
		config, err := cfg.ReadConfig(string(t.Config.SourceConfig))
		if err != nil {
			return nil, err
		}
		if config.Project == "" {
			// fill-in missing project with our project
			config.Project = t.Config.Project
		}
		t.sourceConfig = config
	}
	return t.sourceConfig, nil
}

func (t *Instance) GetOwner() (string, error) {
	return t.Config.DeviceOwner.Substitute(t.Properties)
}

func (t *Instance) NewDeviceMap() (map[string]*srv.Device, error) {
	devices := make(map[string]*srv.Device)

	for deviceName, device := range t.Config.Devices {
		dir, err := t.DeviceDir(deviceName, device)
		if err != nil {
			return nil, err
		}
		devices[deviceName] = &srv.Device{
			Path:     device.Path,
			Source:   dir,
			Readonly: device.Readonly,
		}
	}
	return devices, nil
}

func (instance *Instance) PrintDevices() error {
	writer := &table.FixedWriter{Writer: os.Stdout}
	devices, err := instance.DeviceList()
	if err != nil {
		return err
	}
	var d InstanceDevice
	writer.Columns(
		table.NewColumn("SOURCE", func() interface{} { return d.Source }),
		table.NewColumn("PATH", func() interface{} { return d.Device.Path }),
		table.NewColumn("NAME", func() interface{} { return d.Name }),
		table.NewColumn("DIR", func() interface{} { return d.Device.Dir }),
		table.NewColumn("FILESYSTEM", func() interface{} { return d.Device.Filesystem }),
	)
	for _, d = range devices {
		writer.WriteRow()
	}
	writer.End()
	return nil
}
