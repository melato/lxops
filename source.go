package lxops

import (
	"errors"
	"fmt"
	"strings"
)

type ContainerSource struct {
	Project   string
	Container string
	Snapshot  string
}

type DeviceSource struct {
	Instance *Instance
	Snapshot string
	Clone    bool
}

func (t *DeviceSource) IsDefined() bool {
	return t.Instance != nil
}

func (t *DeviceSource) String() string {
	return fmt.Sprintf("instance=%s snapshot=%s clone=%v", t.Instance.Name, t.Snapshot, t.Clone)
}

func (t *ContainerSource) parse(s string) {
	i := strings.Index(s, "/")
	pc := s
	if i >= 0 {
		t.Snapshot = s[i+1:]
		pc = s[0:i]
	}
	i = strings.Index(pc, "_")
	if i >= 0 {
		t.Project = s[0:i]
		t.Container = s[i+1:]
	} else {
		t.Container = pc
	}
}

func (t *ContainerSource) String() string {
	//return t.Project + "_" + t.Container + "/" + t.Snapshot
	return fmt.Sprintf("project=%s container=%s snapshot=%s", t.Project, t.Container, t.Snapshot)
}

func (t *ContainerSource) IsDefined() bool {
	return t.Container != ""
}

func (t *Instance) newDeviceSource() (*DeviceSource, error) {
	config := t.Config
	source := &DeviceSource{}
	var name string
	if config.DeviceOrigin != "" {
		// check device-origin first
		s, err := config.DeviceOrigin.Substitute(t.Properties)
		if err != nil {
			return nil, err
		}
		parts := strings.Split(s, "@")
		if len(parts) != 2 || len(parts[1]) == 0 {
			return nil, errors.New("missing device origin snapshot: " + s)
		}
		name = parts[0]
		source.Snapshot = parts[1]
		source.Clone = true
	} else if config.DeviceTemplate != "" {
		// device-template is used only if device-origin is missing
		var err error
		name, err = config.DeviceTemplate.Substitute(t.Properties)
		if err != nil {
			return nil, err
		}
	} else {
		return source, nil
	}
	if name == "" && config.SourceConfig == "" {
		return nil, errors.New("missing devices source name")
	}
	sourceConfig, err := t.GetSourceConfig()
	if err != nil {
		return nil, err
	}
	source.Instance, err = newInstance(t.GlobalProperties, sourceConfig, name, false)
	if err != nil {
		return nil, err
	}
	return source, nil
}

func (t *Instance) newContainerSource() (*ContainerSource, error) {
	s, err := t.Config.Origin.Substitute(t.Properties)
	if err != nil {
		return nil, err
	}
	source := &ContainerSource{}
	if s != "" {
		source.parse(s)
		if source.Project == "" {
			sourceConfig, err := t.GetSourceConfig()
			if err != nil {
				return nil, err
			}
			source.Project = sourceConfig.Project
		}
		if source.Project == "" {
			source.Project = t.Config.Project
		}
	}
	return source, nil
}
