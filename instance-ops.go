package lxops

import (
	"fmt"
	"os"

	"melato.org/lxops/util"
	"melato.org/table3"
)

type InstanceOps struct {
	Client ConfigContext `name:"-"`
	ConfigOptions
	Trace bool
	DryRunFlag
}

func (t *InstanceOps) Init() error {
	return t.ConfigOptions.Init()
}

func (t *InstanceOps) Configured() error {
	if t.DryRun {
		t.Trace = true
	}
	t.ConfigOptions.ConfigureProject(t.Client)
	return t.ConfigOptions.Configured()
}

func (t *InstanceOps) Verify(instance *Instance) error {
	fmt.Println(instance.Name)
	return nil
}

// Description prints the description of the instance
func (t *InstanceOps) Description(instance *Instance) error {
	fmt.Println(instance.Config.Description)
	return nil
}

// Properties prints the instance properties
func (t *InstanceOps) Properties(instance *Instance) error {
	fmt.Printf("Effective Properties:\n")
	util.PrintProperties(instance.EffectiveProperties())
	return nil
}

func (t *InstanceOps) Filesystems(instance *Instance) error {
	filesystems, err := instance.FilesystemList()
	if err != nil {
		return err
	}
	writer := &table.FixedWriter{Writer: os.Stdout}
	var fs *InstanceFS
	writer.Columns(
		table.NewColumn("FILESYSTEM", func() interface{} { return fs.Id }),
		table.NewColumn("PATH", func() interface{} { return fs.Path }),
		table.NewColumn("PATTERN", func() interface{} { return fs.Filesystem.Pattern }),
	)
	for _, fs = range filesystems {
		writer.WriteRow()
	}
	writer.End()
	return nil
}

func (t *InstanceOps) Devices(instance *Instance) error {
	return instance.PrintDevices()
}

func (t *InstanceOps) Project(instance *Instance) error {
	fmt.Println(instance.Config.Project)
	return nil
}
