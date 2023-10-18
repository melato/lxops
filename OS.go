package lxops

import (
	"fmt"

	"melato.org/cloudconfig"
)

var OSTypes map[string]cloudconfig.OSType

func init() {
	OSTypes = make(map[string]cloudconfig.OSType)
}

func (t *OS) Type() cloudconfig.OSType {
	osType, exists := OSTypes[t.Name]
	if !exists {
		fmt.Println("Unknown OS type: " + t.Name)
	}
	return osType
}

func (t *OS) IsAlpine() bool {
	return t.Name == "alpine"
}
