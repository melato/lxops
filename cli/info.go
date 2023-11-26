package cli

import (
	"fmt"
	"sort"

	"melato.org/lxops/cfg"
)

type InfoOps struct {
}

func (t *InfoOps) ListOSTypes() {
	types := make([]string, 0, len(cfg.OSTypes))
	for name, _ := range cfg.OSTypes {
		types = append(types, name)
	}
	sort.Strings(types)
	for _, name := range types {
		fmt.Printf("%s\n", name)
	}
}
