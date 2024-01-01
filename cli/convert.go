package cli

import (
	"fmt"
	"os"

	"melato.org/lxops/cfg"
	"melato.org/lxops/yaml"
)

type ConvertOps struct {
	OutputFile string `name:"o" usage:"output file"`
	Force      bool   `name:"f" usage:"force conversion even if existing format is the latest"`
}

func (t *ConvertOps) ConvertFile(inputFile, outputFile string) error {
	data, err := os.ReadFile(inputFile)
	if err != nil {
		return err
	}
	if !t.Force && inputFile == outputFile {
		comment := yaml.FirstLineComment(data)
		if comment == cfg.Comment {
			return nil
		}
	}
	c, err := cfg.Unmarshal(data)
	if err != nil {
		return err
	}
	data, err = cfg.Marshal(c)
	if err != nil {
		return err
	}
	fmt.Printf("%s\n", outputFile)
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
