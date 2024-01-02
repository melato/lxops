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

type RootFS struct {
	OriginFilesystem string
	Instance         string
	Mountpoint       string
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

func NewRootFS(server srv.InstanceServer, instance string) (*RootFS, error) {
	root, err := findRootStoragePool(server, instance)
	if err != nil {
		return nil, err
	}
	if root.Driver != "zfs" {
		return nil, fmt.Errorf("root storage pool is not zfs")
	}
	fs := path.Join(root.Source, "containers", instance)
	var rootfs RootFS
	s := rootfs.newScript()
	rootfs.OriginFilesystem = s.Cmd("zfs", "list", "-H", "-o", "origin", fs).ToString()
	if s.HasError() {
		return nil, s.Error()
	}
	if strings.IndexByte(rootfs.OriginFilesystem, byte('c')) < 0 {
		return nil, fmt.Errorf("origin filesystem is not a snapshot: %s", rootfs.OriginFilesystem)
	}
	rootfs.Instance = instance
	return &rootfs, nil
}

func (t *RootFS) Mount() error {
	if t.Mountpoint != "" {
		return nil // already mounted
	}
	t.Mountpoint = "/tmp/lxops/" + t.Instance
	s := t.newScript()
	s.Run("sudo", "mount", "-t", "zfs", t.OriginFilesystem, t.Mountpoint)
	return s.Error()
}

func (t *RootFS) newScript() *script.Script {
	return &script.Script{Trace: true}
}

func (t *RootFS) IsMounted() bool {
	return t.Mountpoint != ""
}

func (t *RootFS) Unmount() error {
	if !t.IsMounted() {
		return nil
	}
	s := t.newScript()
	s.Run("sudo", "umount", t.Mountpoint)
	t.Mountpoint = ""
	return s.Error()
}

func (t *RootFS) MountedDir(dir string) (string, error) {
	if !t.IsMounted() {
		return "", fmt.Errorf("not mounted")
	}
	fdir := filepath.Join(t.Mountpoint, "rootfs", dir)
	finfo, err := os.Stat(fdir)
	if err != nil {
		return "", err
	}
	if !finfo.IsDir() {
		return "", fmt.Errorf("not a directory: %s", fdir)
	}
	return fdir, nil
}
