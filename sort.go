package lxops

import (
	"sort"
	"strings"
)

type NamedDevice struct {
	Name   string
	Device *Device
}

type DeviceSorter []NamedDevice

func (t DeviceSorter) Len() int {
	return len(t)
}

func (t DeviceSorter) Less(i, j int) bool {
	d1 := t[i]
	d2 := t[j]
	c := strings.Compare(d1.Device.Filesystem, d2.Device.Filesystem)
	if c == 0 {
		return d1.Name < d2.Name
	}
	return c < 0
}

func (t DeviceSorter) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

func SortDevices(devices map[string]*Device) []NamedDevice {
	nd := make([]NamedDevice, 0, len(devices))
	for name, device := range devices {
		nd = append(nd, NamedDevice{name, device})
	}
	sort.Sort(DeviceSorter(nd))
	return nd
}
