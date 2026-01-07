package lxops

import (
	"errors"
	"os/exec"
	"path/filepath"
	"time"

	"melato.org/lxops/cfg"
	"melato.org/script"
)

type Migrate struct {
	PropertyOptions
	FromHost      string `name:"from-host" usage:"optional source host"`
	ToHost        string `name:"to-host" usage:"optional destination host"`
	ConfigFile    string `name:"c" usage:"configFile"`
	FromContainer string `name:"from-container" usage:"source instance"`
	Container     string `name:"container" usage:"destination instance"`
	Snapshot      string `name:"s" usage:"snapshot name"`
	DryRunFlag
	makeSnapshot bool
}

func (t *Migrate) Init() error {
	return t.PropertyOptions.Init()
}

func (t *Migrate) Configured() error {
	if t.ConfigFile == "" {
		return errors.New("missing config file")
	}
	if t.Container == "" {
		return errors.New("missing container")
	}
	if t.FromContainer == "" {
		t.FromContainer = t.Container
	}
	if t.FromHost != "" && t.ToHost != "" {
		return errors.New("cannot use both -from-host and -to-host")
	}
	if t.Snapshot == "" {
		t.Snapshot = time.Now().UTC().Format("20060102150405")
		t.makeSnapshot = true
		if !filepath.IsAbs(t.ConfigFile) {
			return errors.New("config file should be absolute to make a remote snapshot ")
		}
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

func (t *Migrate) GetProperty(name string) (string, bool) {
	value, found := t.GlobalProperties[name]
	return value, found
}

func (t *Migrate) CopyFilesystems() error {
	var config *cfg.Config
	config, err := cfg.ReadConfig(t.ConfigFile, t.GetProperty)
	if err != nil {
		return err
	}
	instance, err := NewInstance(nil, t.GlobalProperties, config, t.Container)
	if err != nil {
		return err
	}
	fromInstance := instance
	if t.FromContainer != t.Container {
		fromInstance, err = NewInstance(nil, t.GlobalProperties, config, t.FromContainer)
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
		s.RunCmd(t.hostCommand(t.FromHost, "lxops", "snapshot", "-s", t.Snapshot, "--name", t.FromContainer, t.ConfigFile))
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
