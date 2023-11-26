package cfg

import (
	"fmt"

	"melato.org/cloudconfig"
)

var OSTypes map[string]cloudconfig.OSType

func init() {
	OSTypes = make(map[string]cloudconfig.OSType)
}

func OSType(name string) (cloudconfig.OSType, error) {
	ostype, exists := OSTypes[name]
	if !exists {
		return nil, fmt.Errorf("unsupported OS type: %s", name)

	}
	return ostype, nil
}
