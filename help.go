package lxops

import (
	"io/fs"

	"melato.org/command"
	"melato.org/lxops/cfg"
	"melato.org/lxops/internal/doc"
)

func helpCommand(fsys fs.FS) *command.SimpleCommand {
	cmd := &command.SimpleCommand{}
	app := &doc.Doc{FS: fsys}
	cmd.Flags(app)
	cmd.Command("config").RunFunc(app.PrintTypeFunc((*cfg.Config)(nil), "Config"))
	cmd.Command("filesystem").RunFunc(app.PrintTypeFunc((*cfg.Filesystem)(nil), "Filesystem"))
	cmd.Command("device").RunFunc(app.PrintTypeFunc((*cfg.Device)(nil), "Device"))
	cmd.Command("pattern").RunFunc(app.PrintTypeFunc((*cfg.Pattern)(nil), "Pattern"))
	cmd.Command("hostpath").RunFunc(app.PrintTypeFunc((*cfg.HostPath)(nil), "HostPath"))
	return cmd
}
