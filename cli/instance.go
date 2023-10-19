package cli

import (
	"fmt"
	"os"
	"sort"

	"melato.org/lxops/srv"
	table "melato.org/table3"
)

// InstanceOps - operations on instances
type InstanceOps struct {
	Client srv.Client `name:"-"`
	server srv.InstanceServer
}

func (t *InstanceOps) Configured() error {
	server, err := t.Client.CurrentInstanceServer()
	if err != nil {
		return err
	}
	t.server = server
	return nil
}

func (t *InstanceOps) Profiles(instance string) error {
	profiles, err := t.server.GetInstanceProfiles(instance)
	if err != nil {
		return err
	}
	for _, profile := range profiles {
		fmt.Println(profile)
	}
	return nil
}

func (t *InstanceOps) Wait(args []string) error {
	for _, instance := range args {
		err := t.server.WaitForNetwork(instance)
		if err != nil {
			return err
		}
	}
	return nil
}

type disk_device struct {
	Name string
	srv.Device
}
type disk_device_sorter []disk_device

func (t disk_device_sorter) Len() int           { return len(t) }
func (t disk_device_sorter) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }
func (t disk_device_sorter) Less(i, j int) bool { return t[i].Source < t[j].Source }

func (t *InstanceOps) Devices(instance string) error {
	devs, err := t.server.GetInstanceDevices(instance)
	if err != nil {
		return err
	}
	writer := &table.FixedWriter{Writer: os.Stdout}

	var devices []disk_device
	for name, d := range devs {
		devices = append(devices, disk_device{Name: name, Device: *d})
	}
	sort.Sort(disk_device_sorter(devices))

	var d disk_device
	writer.Columns(
		table.NewColumn("SOURCE", func() interface{} { return d.Source }),
		table.NewColumn("PATH", func() interface{} { return d.Path }),
		table.NewColumn("NAME", func() interface{} { return d.Name }),
		table.NewColumn("READONLY", func() interface{} { return d.Readonly }),
	)
	for _, d = range devices {
		writer.WriteRow()
	}
	writer.End()
	return nil
}

func (t *InstanceOps) ListHwaddr() error {
	addresses, err := t.server.GetHwaddresses()
	if err != nil {
		return err
	}
	var a srv.Hwaddr
	writer := &table.FixedWriter{Writer: os.Stdout}
	writer.Columns(
		table.NewColumn("HWADDR", func() interface{} { return a.Hwaddr }),
		table.NewColumn("NAME", func() interface{} { return a.Instance }),
	)
	for _, a = range addresses {
		writer.WriteRow()
	}
	writer.End()
	return nil
}

func (t *InstanceOps) ListImages() error {
	list, err := t.server.GetInstanceImages()
	if err != nil {
		return err
	}
	var im srv.InstanceImage
	writer := &table.FixedWriter{Writer: os.Stdout}
	writer.Columns(
		table.NewColumn("IMAGE", func() interface{} {
			return im.Image
		}),
		table.NewColumn("INSTANCE", func() interface{} { return im.Instance }),
	)
	for _, im = range list {
		writer.WriteRow()
	}
	writer.End()
	return nil
}

func (t *InstanceOps) PublishInstance(instance, snapshot, alias string) error {
	return t.server.PublishInstance(instance, snapshot, alias)
}
