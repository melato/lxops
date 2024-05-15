package cli

import (
	"time"
)

// ImageMetadataOps - edit image metadata
type ImageMetadataOps struct {
	Properties ImageMetadataOptions
	File       string `name:"f" usage:"metadata.yaml file`
}

func (t *ImageMetadataOps) Init() error {
	t.Properties.Init()
	return nil
}

func (t *ImageMetadataOps) Update() error {
	var m *ImageMetadata
	var err error
	if t.File != "" {
		m, err = ReadImageMetadata(t.File)
	} else {
		m, err = NewImageMetadata()
	}
	if err != nil {
		return err
	}
	f := m.GetFields()
	t.Properties.Override(f)
	t.Properties.SetImageDescription(f, t.Properties.Variant)
	date := time.Now()
	if f.Serial == "" || f.Serial == "." {
		f.Serial = t.Properties.FormatSerial(date)
	}
	m.SetFields(f)
	m.SetDates(date, t.Properties.ExpiryDays)
	if t.File != "" {
		return m.WriteFile(t.File)
	} else {
		return m.Print()
	}
}
