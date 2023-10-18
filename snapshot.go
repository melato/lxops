package lxops

import (
	"errors"
	"time"

	"melato.org/script"
)

type SnapshotParams struct {
	DryRun    bool   `name:"dry-run" usage:"show the commands to run, but do not change anything"`
	Snapshot  string `name:"s" usage:"short snapshot name"`
	Container bool   `name:"c" usage:"also create container snapshot"`
	Destroy   bool   `name:"d" usage:"destroy snapshots"`
	Recursive bool   `name:"R" usage:"zfs destroy -R: Recursively destroy all dependents, including cloned datasets"`
}

type Snapshot struct {
	Client ConfigContext `name:"-"`
	ConfigOptions
	SnapshotParams
}

func (t *Snapshot) Init() error {
	t.SnapshotParams.Snapshot = time.Now().UTC().Format("20060102150405")
	return t.ConfigOptions.Init()
}

func (t *Snapshot) Configured() error {
	if len(t.SnapshotParams.Snapshot) == 0 {
		return errors.New("empty snapshot name")
	}
	if t.Recursive && !t.Destroy {
		return errors.New("cannot use -R without -d")
	}
	return t.ConfigOptions.Configured(t.Client)
}

func (t *Snapshot) DestroySnapshot(instance *Instance) error {
	filesystems, err := instance.FilesystemList()
	if err != nil {
		return err
	}
	s := &script.Script{Trace: true}
	if t.Recursive {
		roots := InstanceFSList(filesystems).Roots()
		for _, fs := range roots {
			s.Run("sudo", "zfs", "destroy", "-R", fs.Path+"@"+t.Snapshot)
		}
	} else {
		for _, fs := range filesystems {
			s.Run("sudo", "zfs", "destroy", fs.Path+"@"+t.Snapshot)
		}
	}
	return s.Error()
}

func (t *Snapshot) Run(instance *Instance) error {
	if t.Destroy {
		return t.DestroySnapshot(instance)
	} else {
		if t.Container {
			s := &script.Script{Trace: true}
			s.Run("lxc", "snapshot", instance.Container(), t.Snapshot)
			if s.HasError() {
				return s.Error()
			}
		}
		return instance.Snapshot(t.Snapshot)
	}
}
