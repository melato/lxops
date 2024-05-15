package cli

import (
	"os"

	"melato.org/lxops/srv"
)

// ImageExportOps - export images
type ImageExportOps struct {
	Client  srv.Client `name:"-"`
	server  srv.InstanceServer
	Convert ImageConvertOps
	//Squashfs   bool   `name:"squashfs" usage:"export to squashfs"`
}

func (t *ImageExportOps) Init() error {
	return t.Convert.Init()
}

func (t *ImageExportOps) Configured() error {
	server, err := t.Client.CurrentInstanceServer()
	if err != nil {
		return err
	}
	t.server = server
	return t.Convert.Configured()
}

func (t *ImageExportOps) Export(image string) error {
	if t.Convert.Parse != "" {
		err := t.Convert.Properties.ParsePrefixNameTime(image, t.Convert.Parse)
		if err != nil {
			return err
		}
	}
	err := os.MkdirAll(t.Convert.Dir, os.FileMode(0775))
	if err != nil {
		return err
	}

	exportDir, err := mkdir(t.Convert.Dir, "export")
	if err != nil {
		return err
	}

	err = t.server.ExportImage(image, exportDir)
	if err != nil {
		return err
	}

	tarfile, err := findTarfile(exportDir)
	if err != nil {
		return err
	}
	err = t.Convert.ConvertTarfile(tarfile)
	if err != nil {
		return err
	}
	if !t.Convert.Keep {
		err = os.RemoveAll(exportDir)
		if err != nil {
			return err
		}
	}
	return nil
}
