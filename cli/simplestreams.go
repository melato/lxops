package cli

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

type SimplestreamsOps struct {
	Dir string `name:"d" usage:"simplestreams directory"`
}

func (t *SimplestreamsOps) Add(metadataTarball, datafile string) error {
	if t.Dir != "" {
		var err error
		metadataTarball, err = filepath.Abs(metadataTarball)
		if err == nil {
			datafile, err = filepath.Abs(datafile)
		}
		if err != nil {
			return err
		}
	}
	cmd := exec.Command("incus-simplestreams", "add", metadataTarball, datafile)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = t.Dir
	fmt.Printf("%s\n", cmd.String())
	return cmd.Run()
}
