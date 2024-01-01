package lxops

import (
	"io/fs"

	"melato.org/command"
	"melato.org/lxops/cfg"
	"melato.org/lxops/internal/doc"
)

func helpCommand(helpFS fs.FS) *command.SimpleCommand {
	cmd := &command.SimpleCommand{}
	app := &doc.Doc{FS: helpFS}
	cmd.Flags(app)
	cmd.Command("config").RunFunc(func() { app.PrintType((*cfg.Config)(nil), "Config") })
	cmd.Command("filesystem").RunFunc(func() { app.PrintType((*cfg.Filesystem)(nil), "Filesystem") })
	cmd.Command("device").RunFunc(func() { app.PrintType((*cfg.Device)(nil), "Device") })
	cmd.Command("pattern").RunFunc(func() { app.PrintType((*cfg.Pattern)(nil), "Pattern") })
	cmd.Command("hostpath").RunFunc(func() { app.PrintType((*cfg.HostPath)(nil), "HostPath") })

	topics := doc.NewTopics(helpFS, "topics", ".tpl")
	cmd.Command("topics").RunFunc(topics.Print)
	return cmd
}
