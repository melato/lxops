package cli

import (
	"fmt"
	"os"

	"melato.org/lxops/cfg"
)

type ConvertOps struct {
	OutputFile string `name:"o" usage:"output file"`
}

func (t *ConvertOps) ConvertFile(inputFile, outputFile string) error {
	c, err := cfg.ReadConfigYaml(inputFile)
	if err != nil {
		return err
	}
	data, err := cfg.Marshal(c)
	if err != nil {
		return err
	}
	return os.WriteFile(outputFile, data, os.FileMode(0664))
}

func (t *ConvertOps) Convert(args ...string) error {
	if t.OutputFile != "" {
		if len(args) != 1 {
			return fmt.Errorf("-o requires one input file only")
		}
		return t.ConvertFile(args[0], t.OutputFile)
	} else {
		for _, file := range args {
			err := t.ConvertFile(file, file)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
