package cli

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Exec struct {
	Verbose bool `name:"verbose" usage:"print output of commands"`
}

func (t *Exec) Run(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	if t.Verbose {
		cmd.Stdout = os.Stdout
	}
	cmd.Stderr = os.Stderr
	fmt.Printf("%s\n", cmd.String())
	return cmd.Run()
}

func (t *Exec) ExtractTar(tarfile, unpackDir string) error {
	tarFlags := "xf"
	if t.Verbose {
		tarFlags += "v"
	}
	return t.Run("sudo", "tar", tarFlags, tarfile, "-C", unpackDir)
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

const TarSuffix = ".tar.gz"

func verifyTarfile(path string) error {
	if !strings.HasSuffix(path, TarSuffix) {
		return fmt.Errorf("%s: does not have .tar.gz suffix", path)
	}
	return nil
}

func findTarfile(dir string) (string, error) {
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
	name := entries[0].Name()
	err = verifyTarfile(name)
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, name), nil
}

func mkdir(parent, dir string) (string, error) {
	path := filepath.Join(parent, dir)
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
