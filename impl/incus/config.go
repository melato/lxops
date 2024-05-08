package lxops_incus

import (
	"os"
	"path/filepath"

	config "github.com/lxc/incus/v6/shared/cliconfig"
	"melato.org/lxops/yaml"
)

type Config struct {
	currentProject string
}

func ConfigDir() (string, error) {
	configDir := os.Getenv("INCUS_CONF")
	if configDir != "" {
		return configDir, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	configDir = filepath.Join(home, ".config", "incus")
	if _, err = os.Stat(configDir); err == nil {
		return configDir, nil
	}
	return "", err
}

func (t *Config) getCurrentProject() (string, error) {
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

func (t *Config) CurrentProject() string {
	if t.currentProject == "" {
		project, err := t.getCurrentProject()
		if err != nil || project == "" {
			project = "default"
		}
		t.currentProject = project
	}
	return t.currentProject
}
