package cli

import (
	"fmt"
	"syscall"
	"time"

	"melato.org/lxops/srv"
	"melato.org/lxops/yaml"
)

// ImageMetadataOps - edit image metadata
type ImageMetadataOps struct {
	File       string `name:"f" usage:"metadata.yaml file`
	Name       string `name:"name" usage:"a simple name for the image"`
	Release    string `name:"release" usage:"overrides image.release"`
	OS         string `name:"os" usage:"overrides image.os"`
	ExpireDays int    `name:"expire-days" usage:"expiry_date as number of days from creation_date"`
}

func SetImageProperties(im *srv.ImageFields, name string, tm time.Time) {
	serial := tm.UTC().Format("20060102_15:04")
	im.Serial = serial
	im.Variant = name
	im.Name = fmt.Sprintf("%s-%s-%s-%s-%s", im.OS, im.Release, im.Architecture, im.Variant, im.Serial)
	im.Description = fmt.Sprintf("%s %s %s (%s)", capFirst(name), im.Release, im.Architecture, im.Serial)
	// Release
	// Architecture
	// OS

	/*
	  image.description: Alpinelinux 3.19 x86_64 (20240129_1300)
	  image.name: alpinelinux-3.19-x86_64-default-20240129_1300
	  image.os: alpinelinux
	  image.release: "3.19"
	  image.serial: "20240129_1300"
	  image.variant: default
	*/
}

func int65toString(a *[65]int8) string {
	b := make([]byte, 0, 65)
	for _, c := range a {
		if c == 0 {
			break
		}
		b = append(b, byte(c))
	}
	return string(b)
}

func getSystemArchitecture() (string, error) {
	var u syscall.Utsname
	err := syscall.Uname(&u)
	if err != nil {
		return "", err
	}
	return int65toString(&u.Machine), nil
}

func UpdateImageProperties(im *srv.ImageFields, properties map[any]any) {
	update := func(name, value string) {
		if value != "" {
			properties[name] = value
		}
	}
	update("architecture", im.Architecture)
	update("description", im.Description)
	update("name", im.Name)
	update("os", im.OS)
	update("release", im.Release)
	update("serial", im.Serial)
	update("variant", im.Variant)
}

func (t *ImageMetadataOps) Init() error {
	t.ExpireDays = 30
	return nil
}

func (t *ImageMetadataOps) Update() error {
	if t.Name == "" {
		return fmt.Errorf("missing name")
	}
	var v map[any]any
	var properties map[any]any
	var im srv.ImageFields
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
		im.Architecture, err = getSystemArchitecture()
		if err != nil {
			return err
		}
	}
	if im.Release == "" {
		im.Release = t.Release
	}
	if im.OS == "" {
		im.OS = t.OS
	}
	date := time.Now()
	timestamp := date.Unix()
	v["creation_date"] = timestamp
	v["expiry_date"] = timestamp + int64(t.ExpireDays)*24*3600
	SetImageProperties(&im, t.Name, date)
	UpdateImageProperties(&im, properties)
	if t.File != "" {
		return yaml.WriteFile(v, t.File)
	} else {
		return yaml.Print(v)
	}
}
