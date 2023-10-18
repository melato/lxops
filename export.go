package lxops

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"melato.org/lxops/srv"
)

var TraceExport bool

type ExportOps struct {
	Client srv.Client `name:"-"`
	ConfigOptions
	Dir      string `name:"d" usage:"import/export directory"`
	Snapshot string `name:"snapshot" usage:"short name of snapshot to export"`
	Image    bool   `name:"image" usage:"export/import lxc image too -- experimental"`
	DryRun   bool
}

func (t *ExportOps) Run(name string, arg ...string) error {
	cmd := exec.Command(name, arg...)
	if TraceExport {
		fmt.Printf("%s\n", cmd.String())
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (t *ExportOps) Export(configFile string) error {
	TraceExport = true
	instance, err := t.Instance(configFile)
	if err != nil {
		return err
	}
	dir := t.Dir
	filesystems, err := instance.FilesystemList()
	if err != nil {
		return err
	}
	if t.Image {
		server, err := t.Client.CurrentInstanceServer()
		if err != nil {
			return err
		}
		err = server.ExportImage(instance.Name, filepath.Join(dir, "image"))
		if err != nil {
			return err
		}
	}
	err = t.Run("sudo", "mkdir", "-p", dir)
	if err != nil {
		return err
	}

	var mntDir string
	if t.Snapshot != "" {
		mntDir = filepath.Join(dir, "mnt")
		err := t.Run("sudo", "mkdir", "-p", mntDir)
		if err != nil {
			return err
		}
		defer t.Run("sudo", "rmdir", mntDir)
	}

	for _, fs := range filesystems {
		if fs.Filesystem.Transient {
			continue
		}
		tarFile := filepath.Join(dir, fs.Id+".tar.gz")
		var err error
		if t.Snapshot == "" || fs.IsDir() {
			err = t.Run("sudo", "tar", "cfz", tarFile, "-C", fs.Dir(), ".")
		} else {
			err = t.Run("sudo", "mount", "-t", "zfs", "-o", "ro", fs.Path+"@"+t.Snapshot, mntDir)
			if err != nil {
				return err
			}
			err = t.Run("sudo", "tar", "cfz", tarFile, "-C", mntDir, ".")
			err2 := t.Run("sudo", "umount", mntDir)
			if err == nil {
				err = err2
			}
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *ExportOps) Import(configFile string) error {
	instance, err := t.Instance(configFile)
	if err != nil {
		return err
	}
	dir := t.Dir
	filesystems, err := instance.FilesystemList()
	if err != nil {
		return err
	}
	if t.Image {
		err := t.Run("lxc", "image", "import", filepath.Join(dir, instance.Name+".tar.gz"), "--alias="+instance.Name)
		if err != nil {
			return err
		}
	}
	dev, err := NewDeviceConfigurer(instance)
	if err != nil {
		return err
	}
	dev.NoRsync = true
	dev.Trace = TraceExport
	dev.DryRun = t.DryRun
	err = dev.ConfigureDevices(instance)
	if err != nil {
		return err
	}

	for _, fs := range filesystems {
		if fs.Filesystem.Transient {
			continue
		}
		tarFile := filepath.Join(dir, fs.Id+".tar.gz")
		err = t.Run("sudo", "tar", "xfz", tarFile, "-C", fs.Dir(), ".")
		if err != nil {
			return err
		}
	}
	return nil
}
