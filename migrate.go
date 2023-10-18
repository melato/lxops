package lxops

import (
	"errors"
	"os/exec"
	"path/filepath"
	"time"

	"melato.org/script"
)

type Migrate struct {
	PropertyOptions
	FromHost      string
	ToHost        string
	ConfigFile    string `name:"c" usage:"configFile"`
	FromContainer string
	Container     string
	Snapshot      string `name:"s" usage:"snapshot name"`
	DryRun        bool   `name:"dry-run" usage:"show the commands to run, but do not change anything"`
	makeSnapshot  bool
}

func (t *Migrate) Init() error {
	return t.PropertyOptions.Init()
}

func (t *Migrate) Configured() error {
	if len(t.Snapshot) == 0 {
		t.Snapshot = time.Now().UTC().Format("20060102150405")
		t.makeSnapshot = true
	}
	if (t.FromHost == "") == (t.ToHost == "") {
		return errors.New("need to provide either -from-host or -to-host")
	}
	if t.ConfigFile == "" {
		return errors.New("missing config file")
	}
	if t.Container == "" {
		return errors.New("missing container")
	}
	if t.FromContainer == "" {
		t.FromContainer = t.Container
	}
	if !filepath.IsAbs(t.ConfigFile) {
		return errors.New("config file should be absolute")
	}
	return t.PropertyOptions.Configured()
}

func (t *Migrate) hostCommand(host, command string, args ...string) *exec.Cmd {
	if host != "" {
		return exec.Command("ssh", append([]string{host, command}, args...)...)
	} else {
		return exec.Command(command, args...)
	}
}

func (t *Migrate) CopyFilesystems() error {
	var config *Config
	config, err := ReadConfig(t.ConfigFile)
	if err != nil {
		return err
	}
	instance, err := NewInstance(t.GlobalProperties, config, t.Container)
	if err != nil {
		return err
	}
	fromInstance := instance
	if t.FromContainer != t.Container {
		fromInstance, err = NewInstance(t.GlobalProperties, config, t.FromContainer)
		if err != nil {
			return err
		}
	}

	filesystems, err := instance.FilesystemList()
	if err != nil {
		return err
	}
	fromFilesystems, err := fromInstance.Filesystems()
	if err != nil {
		return err
	}
	s := script.Script{Trace: true, DryRun: t.DryRun}
	if t.makeSnapshot {
		s.RunCmd(t.hostCommand(t.FromHost, "lxdops", "snapshot", "-s", t.Snapshot, "--name", t.FromContainer, t.ConfigFile))
	}
	for _, fs := range filesystems {
		if fs.IsZfs() && !fs.Filesystem.Transient {
			fromFS, ok := fromFilesystems[fs.Id]
			if !ok {
				continue
			}
			send := t.hostCommand(t.FromHost, "sudo", "zfs", "send", fromFS.Path+"@"+t.Snapshot)
			receive := t.hostCommand(t.ToHost, "sudo", "zfs", "receive", fs.Path)
			s.RunCmd(send, receive)
		}
	}
	return s.Error()
}
