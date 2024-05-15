package cli

import (
	"fmt"
	"os"
	"os/exec"
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
	Verbose    bool       `name:"verbose" usage:"print output of commands"`
	Properties ImageMetadataOptions
	Parse      bool `name:"parse" usage:"derive metadata properties from image name"`
	server     srv.InstanceServer
}

func (t *ImageExportOps) Init() error {
	t.Dir = "."
	t.Properties.Init()
	return nil
}

func (t *ImageExportOps) Configured() error {
	server, err := t.Client.CurrentInstanceServer()
	if err != nil {
		return err
	}
	t.server = server
	return nil
}

func (t *ImageExportOps) findTarfile(dir string) (string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return "", err
	}
	if len(entries) != 1 {
		for _, e := range entries {
			name := e.Name()
			fmt.Printf("%s\n", name)
		}
		return "", fmt.Errorf("there are multiple files")
	}
	return filepath.Join(dir, entries[0].Name()), nil
}

func (t *ImageExportOps) exec(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	if t.Verbose {
		cmd.Stdout = os.Stdout
	}
	cmd.Stderr = os.Stderr
	fmt.Printf("%s\n", cmd.String())
	return cmd.Run()
}

func checkFileNotExist(path string) error {
	_, err := os.Stat(path)
	if err == nil {
		return fmt.Errorf("file exists: %s\n", path)
	}
	if os.IsNotExist(err) {
		return nil
	}
	return err
}

func checkFilesNotExist(paths ...string) error {
	for _, path := range paths {
		err := checkFileNotExist(path)
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *ImageExportOps) mkdir(dir string) (string, error) {
	path := filepath.Join(t.Dir, dir)
	err := checkFileNotExist(path)
	if err != nil {
		return "", err
	}
	err = os.Mkdir(path, os.FileMode(0775))
	if err != nil {
		return "", err
	}
	return path, nil
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
	err := os.MkdirAll(t.Dir, os.FileMode(0775))
	if err != nil {
		return err
	}
	var rootfsFile string
	var metadataTarFile string
	var unpackDir string

	exportDir := t.Dir
	if t.Squashfs {
		exportDir, err = t.mkdir("export")
		if err != nil {
			return err
		}
		unpackDir, err = t.mkdir("unpack")
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
		tarfile, err := t.findTarfile(exportDir)
		if err != nil {
			return err
		}
		tarFlags := "xf"
		if t.Verbose {
			tarFlags += "v"
		}
		err = t.exec("sudo", "tar", tarFlags, tarfile, "-C", unpackDir)
		if err != nil {
			return err
		}
		err = t.exec("sudo", "mksquashfs", unpackDir,
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
		err = t.exec("sudo", "chown", owner, rootfsFile)
		if err != nil {
			return err
		}

		metadataFile := filepath.Join(unpackDir, "metadata.yaml")
		err = t.exec("sudo", "chown", owner, metadataFile)
		if err != nil {
			return err
		}
		err = t.updateMetadata(metadataFile)
		if err != nil {
			return err
		}
		err = t.exec("tar", "Jcf",
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
			err = t.exec("sudo", "rm", "-rf", unpackDir)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
