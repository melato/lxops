package cli

import (
	"os"
	"os/user"
	"path/filepath"

	"melato.org/lxops/srv"
)

// ImageExportOps - export images
type ImageExportOps struct {
	Client     srv.Client `name:"-"`
	Squashfs   bool       `name:"squashfs" usage:"export to squashfs"`
	Dir        string     `name:"d" usage:"export directory"`
	Keep       bool       `name:"keep" usage:"do not delete intermediate directories"`
	Properties ImageMetadataOptions
	Exec       Exec
	Parse      bool `name:"parse" usage:"derive metadata properties from image name"`
	server     srv.InstanceServer
}

func (t *ImageExportOps) Init() error {
	t.Dir = "."
	return t.Properties.Init()
}

func (t *ImageExportOps) Configured() error {
	server, err := t.Client.CurrentInstanceServer()
	if err != nil {
		return err
	}
	t.server = server
	return t.Properties.Configured()
}

func (t *ImageExportOps) updateMetadata(file string) error {
	m, err := ReadImageMetadata(file)
	if err != nil {
		return err
	}
	t.Properties.Apply(m)
	return m.WriteFile(file)
}

func (t *ImageExportOps) Export(image string) error {
	if t.Parse {
		err := t.Properties.ParsePrefixNameDateTime(image)
		if err != nil {
			return err
		}
	}
	err := os.MkdirAll(t.Dir, os.FileMode(0775))
	if err != nil {
		return err
	}
	var rootfsFile string
	var metadataTarFile string
	var unpackDir string

	exportDir := t.Dir
	if t.Squashfs {
		exportDir, err = mkdir(t.Dir, "export")
		if err != nil {
			return err
		}
		unpackDir, err = mkdir(t.Dir, "unpack")
		if err != nil {
			return err
		}
		rootfsFile = filepath.Join(t.Dir, "rootfs.squashfs")
		metadataTarFile = filepath.Join(t.Dir, "metadata.tar.xz")
		err = checkFilesNotExist(rootfsFile, metadataTarFile)
		if err != nil {
			return err
		}
	}

	err = t.server.ExportImage(image, exportDir)
	if err != nil {
		return err
	}
	if t.Squashfs {
		tarfile, err := findTarfile(exportDir)
		if err != nil {
			return err
		}
		err = t.Exec.ExtractTar(tarfile, unpackDir)
		if err != nil {
			return err
		}
		err = t.Exec.Run("sudo", "mksquashfs", unpackDir,
			rootfsFile,
			"-noappend", "-comp", "xz", "-b", "1M")
		if err != nil {
			return err
		}
		user, err := user.Current()
		if err != nil {
			return err
		}
		owner := user.Uid + ":" + user.Gid
		err = t.Exec.Run("sudo", "chown", owner, rootfsFile)
		if err != nil {
			return err
		}

		metadataFile := filepath.Join(unpackDir, "metadata.yaml")
		err = t.Exec.Run("sudo", "chown", owner, metadataFile)
		if err != nil {
			return err
		}
		err = t.updateMetadata(metadataFile)
		if err != nil {
			return err
		}
		err = t.Exec.Run("tar", "Jcf",
			metadataTarFile,
			"-C", unpackDir,
			"metadata.yaml",
			"templates/")
		if err != nil {
			return err
		}
		if !t.Keep {
			err = os.RemoveAll(exportDir)
			if err != nil {
				return err
			}
			err = t.Exec.Run("sudo", "rm", "-rf", unpackDir)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
