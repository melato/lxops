package lxops

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"melato.org/lxops/srv"
	"melato.org/script"
)

// RootFS mounts the image filesystem of an instance to a local directory
type RootFS struct {
	server           srv.InstanceServer
	instance         string
	originFilesystem string
	mountpoint       string
}

func findRootStoragePool(server srv.InstanceServer, instance string) (*srv.StoragePool, error) {
	profiles, err := server.GetInstanceProfiles(instance)
	if err != nil {
		return nil, err
	}
	for _, profile := range profiles {
		devices, err := server.GetProfileDevices(profile)
		if err != nil {
			return nil, err
		}
		for _, device := range devices {
			if device.Path == "/" {
				return server.GetStoragePool(device.Pool)
			}
		}
	}
	return nil, fmt.Errorf("missing root device")
}

func (t *RootFS) findOriginFilesystem() error {
	root, err := findRootStoragePool(t.server, t.instance)
	if err != nil {
		return err
	}
	if root.Driver != "zfs" {
		return fmt.Errorf("root storage pool is not zfs")
	}
	fs := path.Join(root.Source, "containers", t.instance)
	s := t.newScript()
	t.originFilesystem = s.Cmd("zfs", "list", "-H", "-o", "origin", fs).ToString()
	if s.HasError() {
		return s.Error()
	}
	if strings.IndexByte(t.originFilesystem, byte('c')) < 0 {
		return fmt.Errorf("origin filesystem is not a snapshot: %s", t.originFilesystem)
	}
	return nil
}

func NewRootFS(server srv.InstanceServer, instance string) *RootFS {
	return &RootFS{server: server, instance: instance}
}

func (t *RootFS) Mount() error {
	if t.mountpoint != "" {
		return nil // already mounted
	}
	if t.originFilesystem == "" {
		err := t.findOriginFilesystem()
		if err != nil {
			return err
		}
	}
	t.mountpoint = "/tmp/lxops/" + t.instance
	err := os.MkdirAll(t.mountpoint, os.FileMode(0777))
	if err != nil {
		return err
	}
	s := t.newScript()
	s.Run("sudo", "mount", "-t", "zfs", t.originFilesystem, t.mountpoint)
	return s.Error()
}

func (t *RootFS) newScript() *script.Script {
	return &script.Script{Trace: true}
}

func (t *RootFS) IsMounted() bool {
	return t.mountpoint != ""
}

func (t *RootFS) Unmount() error {
	if !t.IsMounted() {
		return nil
	}
	s := t.newScript()
	s.Run("sudo", "umount", t.mountpoint)
	err := os.Remove(t.mountpoint)
	if err != nil {
		return err
	}
	t.mountpoint = ""
	return s.Error()
}

func (t *RootFS) MountedDir(dir string) (string, error) {
	if !t.IsMounted() {
		return "", fmt.Errorf("not mounted")
	}
	fdir := filepath.Join(t.mountpoint, "rootfs", dir)
	finfo, err := os.Stat(fdir)
	if err != nil {
		return "", err
	}
	if !finfo.IsDir() {
		return "", fmt.Errorf("not a directory: %s", fdir)
	}
	return fdir, nil
}
