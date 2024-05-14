package cli

import (
	"fmt"
	"time"

	"melato.org/lxops/srv"
	"melato.org/lxops/yaml"
)

// ImageMetadataOps - edit image metadata
type ImageMetadataOps struct {
	File       string `name:"f" usage:"metadata.yaml file`
	Variant    string `name:"variant" usage:"override variant property"`
	Release    string `name:"release" usage:"override release property"`
	OS         string `name:"os" usage:"override os property"`
	Serial     string `name:"serial" usage:"optional serial property"`
	ExpireDays int    `name:"expire-days" usage:"expiry_date as number of days from creation_date"`
}

func (t *ImageMetadataOps) Init() error {
	t.ExpireDays = 30
	return nil
}

func (t *ImageMetadataOps) Update() error {
	var v map[any]any
	var properties map[any]any
	var im srv.ImageFields
	var m ImageMetadata
	if t.File != "" {
		err := yaml.ReadFile(t.File, &v)
		if err != nil {
			return err
		}
		vProperties := v["properties"]
		if vProperties == nil {
			return fmt.Errorf("missing properties")
		}
		var ok bool
		properties, ok = vProperties.(map[any]any)
		if !ok {
			return fmt.Errorf("properties: %T. expected: map[any]any", vProperties)
		}
		getValue := func(name string) string {
			v, ok := properties[name]
			if ok {
				s, ok := v.(string)
				if ok {
					return s
				}
			}
			return ""
		}
		im.Architecture = getValue("architecture")
		im.Description = getValue("description")
		im.Name = getValue("name")
		im.OS = getValue("os")
		im.Release = getValue("release")
		im.Serial = getValue("serial")
		im.Variant = getValue("variant")
	} else {
		v = make(map[any]any)
		properties = make(map[any]any)
		v["properties"] = properties
	}
	if im.Architecture == "" {
		var err error
		im.Architecture, err = m.getSystemArchitecture()
		if err != nil {
			return err
		}
	}
	if t.Release != "" {
		im.Release = t.Release
	}
	if t.OS != "" {
		im.OS = t.OS
	}
	if t.Variant != "" {
		im.Variant = t.Variant
	}
	if t.Serial != "" {
		im.Serial = t.Serial
	}
	date := time.Now()
	timestamp := date.Unix()
	v["creation_date"] = timestamp
	v["expiry_date"] = timestamp + int64(t.ExpireDays)*24*3600
	if im.Serial == "" {
		im.Serial = m.FormatSerial(date)
	}
	m.SetImageDescription(&im, t.Variant)
	m.UpdateImageProperties(&im, properties)
	if t.File != "" {
		return yaml.WriteFile(v, t.File)
	} else {
		return yaml.Print(v)
	}
}
