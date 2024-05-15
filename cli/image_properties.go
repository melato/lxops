package cli

import (
	"fmt"
	"syscall"
	"time"

	"melato.org/lxops/srv"
)

// ImageBaseProperties - independent image properties
type ImageBaseProperties struct {
	Variant string `name:"variant" usage:"override variant property"`
	Release string `name:"release" usage:"override release property"`
	OS      string `name:"os" usage:"override os property"`
	Serial  string `name:"serial" usage:"optional serial property"`
}

func (t *ImageBaseProperties) FormatSerial(tm time.Time) string {
	return tm.UTC().Format("20060102_15:04")
}

// SetImageDescription - generate name, description from other fields
func (t *ImageBaseProperties) SetImageDescription(im *srv.ImageFields, name string) {
	im.Name = fmt.Sprintf("%s-%s-%s-%s-%s", im.OS, im.Release, im.Architecture, im.Variant, im.Serial)
	im.Description = fmt.Sprintf("%s %s %s (%s)", name, im.Release, im.Architecture, im.Serial)

	/* Typical image properties
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

// UpdateImageProperties - copy ImageFields to properties
func (t *ImageBaseProperties) UpdateImageProperties(im *srv.ImageFields, properties map[any]any) {
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

func (t *ImageBaseProperties) Override(im *srv.ImageFields) {
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
}
