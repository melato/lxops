package main

import (
	_ "embed"
	"fmt"

	"melato.org/command"
	"melato.org/command/usage"
	"melato.org/lxops"
	"melato.org/lxops_incus"
)

//go:embed usage.yaml
var usageData []byte

// set with -ldflags "-X 'main.version=...'"
var version = "dev"

func main() {
	lxops.InitOSTypes()
	client := &lxops_incus.Client{}
	lxops.InitConfigTypes()
	cmd := lxops.RootCommand(client)
	cmd.Command("version").NoConfig().RunMethod(func() {
		fmt.Printf("lxops for %s, %s\n", client.ServerType(), version)
	})
	usage.Apply(cmd, usageData)
	command.Main(cmd)
}
