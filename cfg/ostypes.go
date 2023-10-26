package cfg

import (
	"fmt"

	"melato.org/cloudconfig"
)

var OSTypes map[string]cloudconfig.OSType

func init() {
	OSTypes = make(map[string]cloudconfig.OSType)
}

func OSType(ostype string) cloudconfig.OSType {
	os, exists := OSTypes[ostype]
	if !exists {
		fmt.Println("Unknown OS type: " + ostype)
	}
	return os
}
