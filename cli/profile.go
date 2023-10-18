package cli

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"melato.org/lxops/srv"
)

type ProfileOps struct {
	Client srv.Client `name:"-"`
	Dir    string     `name:"d" usage:"export directory"`
}

func (t *ProfileOps) ExportProfile(server srv.InstanceServer, name string) error {
	data, err := server.ExportProfile(name)
	if err != nil {
		return err
	}

	file := path.Join(t.Dir, name)
	return os.WriteFile(file, []byte(data), 0644)
}

func (t *ProfileOps) Export(profiles ...string) error {
	server, err := t.Client.CurrentInstanceServer()
	if err != nil {
		return err
	}
	for _, profile := range profiles {
		err = t.ExportProfile(server, profile)
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *ProfileOps) ImportProfile(server srv.InstanceServer, file string, existingProfiles map[string]bool) error {
	data, err := os.ReadFile(file)
	if err != nil {
		return err
	}
	name := filepath.Base(file)
	return server.ImportProfile(name, data)
}

func (t *ProfileOps) Import(files []string) error {
	server, err := t.Client.CurrentInstanceServer()
	if err != nil {
		return err
	}
	profiles := make(map[string]bool)
	names, err := server.GetProfileNames()
	if err != nil {
		return err
	}
	for _, name := range names {
		profiles[name] = true
	}

	for _, file := range files {
		err := t.ImportProfile(server, file, profiles)
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *ProfileOps) ProfileExists(profile string) error {
	server, err := t.Client.CurrentInstanceServer()
	if err != nil {
		return err
	}
	exists, err := server.ProfileExists(profile)
	if exists {
		fmt.Println(profile)
	}
	return nil
}
