package cli

import (
	"melato.org/lxops/srv"
)

type InstanceImageSorter []srv.InstanceImage

func (t InstanceImageSorter) Len() int      { return len(t) }
func (t InstanceImageSorter) Swap(i, j int) { t[i], t[j] = t[j], t[i] }
func (t InstanceImageSorter) Less(i, j int) bool {
	if t[i].Image != t[j].Image {
		return t[i].Image < t[j].Image
	} else {
		return t[i].Instance < t[j].Instance
	}
}
