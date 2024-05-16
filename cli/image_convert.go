package cli

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

// ImageConvert - convert image file format
type ImageConvertOps struct {
	Dir        string `name:"d" usage:"output and staging directory"`
	Keep       bool   `name:"keep" usage:"do not delete intermediate directories"`
	Properties ImageMetadataOptions
	Exec       Exec
	Parse      string `name:"parse" usage:"extract image metadata from its name"`

	SimplestreamsDir string `name:"simplestreams" usage:"simplestreams directory to add the image to"`
}

func (t *ImageConvertOps) Init() error {
	t.Dir = "."
	return t.Properties.Init()
}

func (t *ImageConvertOps) Configured() error {
	return t.Properties.Configured()
}

func (t *ImageConvertOps) updateMetadata(file string) error {
	fmt.Printf("update %s\n", file)
	m, err := ReadImageMetadata(file)
	if err != nil {
		return err
	}
	t.Properties.Apply(m)
	return m.WriteFile(file)
}

func (t *ImageConvertOps) ConvertTarfile(tarfile string) error {
	err := os.MkdirAll(t.Dir, os.FileMode(0775))
	if err != nil {
		return err
	}
	unpackDir, err := mkdir(t.Dir, "unpack")
	if err != nil {
		return err
	}
	rootfsFile := filepath.Join(t.Dir, "rootfs.squashfs")
	metadataTarFile := filepath.Join(t.Dir, "metadata.tar.xz")
	err = checkFilesNotExist(rootfsFile, metadataTarFile)
	if err != nil {
		return err
	}

	err = t.Exec.ExtractTar(tarfile, unpackDir)
	if err != nil {
		return err
	}
	err = t.Exec.Run("sudo", "mksquashfs", filepath.Join(unpackDir, "rootfs"),
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

	if t.Properties.HasChanges() {
		metadataFile := filepath.Join(unpackDir, "metadata.yaml")
		err = t.Exec.Run("sudo", "chown", owner, metadataFile)
		if err != nil {
			return err
		}
		err = t.updateMetadata(metadataFile)
		if err != nil {
			return err
		}
	}

	err = t.Exec.Run("tar", "Jcf",
		metadataTarFile,
		"-C", unpackDir,
		"metadata.yaml",
		"templates/")
	if err != nil {
		return err
	}
	if t.SimplestreamsDir != "" {
		op := SimplestreamsOps{Dir: t.SimplestreamsDir}
		err := op.Add(metadataTarFile, rootfsFile)
		if err != nil {
			return err
		}
	}
	if !t.Keep {
		err = t.Exec.Run("sudo", "rm", "-rf", unpackDir)
		if err != nil {
			return err
		}
		if t.SimplestreamsDir != "" {
			err = os.Remove(metadataTarFile)
			if err == nil {
				err = os.Remove(rootfsFile)
			}
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (t *ImageConvertOps) ConvertPath(path string) error {
	f, err := os.Stat(path)
	if err != nil {
		return err
	}
	if f.IsDir() {
		tarfile, err := findTarfile(path)
		if err != nil {
			return err
		}
		if t.Parse != "" {
			name := filepath.Base(path)
			err := t.Properties.ParsePrefixNameTime(name, t.Parse)
			if err != nil {
				return err
			}
		}
		return t.ConvertTarfile(tarfile)
	} else {
		err := verifyTarfile(path)
		if err != nil {
			return err
		}
		if t.Parse != "" {
			name := strings.TrimSuffix(filepath.Base(path), TarSuffix)
			err := t.Properties.ParsePrefixNameTime(name, t.Parse)
			if err != nil {
				return err
			}
		}
		return t.ConvertTarfile(path)
	}
}

func (t *ImageConvertOps) Convert(paths ...string) error {
	for _, path := range paths {
		fmt.Printf("converting %s\n", path)
		err := t.ConvertPath(path)
		if err != nil {
			return err
		}
	}
	return nil
}
