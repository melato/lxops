package lxdutil

import (
	"os"
	"path/filepath"

	"github.com/canonical/lxd/lxc/config"
	"melato.org/lxops/yaml"
)

type LxcConfig struct {
	currentProject string
}

func ConfigDir() (string, error) {
	configDir := os.Getenv("LXD_CONF")
	if configDir != "" {
		return configDir, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	configDir = filepath.Join(home, "snap", "lxd", "common", "config")
	if _, err = os.Stat(configDir); err == nil {
		return configDir, nil
	}
	configDir = filepath.Join(home, "snap", "lxd", "current", ".config", "lxc")
	if _, err = os.Stat(configDir); err == nil {
		return configDir, nil
	}
	configDir = filepath.Join(home, ".config", "lxc")
	if _, err = os.Stat(configDir); err == nil {
		return configDir, nil
	}
	return "", err
}

func (t *LxcConfig) getCurrentProject() (string, error) {
	configDir, err := ConfigDir()
	if err != nil {
		return "", err
	}

	var cfg config.Config
	err = yaml.ReadFile(filepath.Join(configDir, "config.yml"), &cfg)
	if err != nil {
		return "", err
	}
	local, found := cfg.Remotes["local"]
	if found {
		return local.Project, nil
	}
	return "", nil
}

func (t *LxcConfig) CurrentProject() string {
	if t.currentProject == "" {
		project, err := t.getCurrentProject()
		if err != nil || project == "" {
			project = "default"
		}
		t.currentProject = project
	}
	return t.currentProject
}
