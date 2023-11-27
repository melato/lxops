package lxops

import (
	"errors"

	"melato.org/script"
)

type Rollback struct {
	Client ConfigContext `name:"-"`
	ConfigOptions
	Snapshot  string `name:"s" usage:"short snapshot name"`
	Container bool   `name:"c" usage:"also restore container snapshot"`
	DryRunFlag
}

func (t *Rollback) Init() error {
	return t.ConfigOptions.Init()
}

func (t *Rollback) Configured() error {
	if t.Snapshot == "" {
		return errors.New("empty snapshot name")
	}
	t.ConfigOptions.ConfigureProject(t.Client)
	return t.ConfigOptions.Configured()
}

func (t *Rollback) Run(instance *Instance) error {
	err := instance.Rollback(t.Snapshot)
	if err != nil {
		return err
	}
	if t.Container {
		s := &script.Script{Trace: true}
		s.Run("lxc", "restore", instance.Container(), t.Snapshot)
		if s.HasError() {
			return s.Error()
		}
	}
	return nil
}
