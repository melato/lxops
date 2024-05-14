package cli

import (
	"fmt"
	"syscall"
	"time"

	"melato.org/lxops/srv"
)

type ImageMetadata struct {
}

func (t *ImageMetadata) FormatSerial(tm time.Time) string {
	return tm.UTC().Format("20060102_15:04")
}

// SetImageDescription - generate name, description from other fields
func (t *ImageMetadata) SetImageDescription(im *srv.ImageFields, name string) {
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

func (t *ImageMetadata) int65toString(a *[65]int8) string {
	b := make([]byte, 0, 65)
	for _, c := range a {
		if c == 0 {
			break
		}
		b = append(b, byte(c))
	}
	return string(b)
}

func (t *ImageMetadata) getSystemArchitecture() (string, error) {
	var u syscall.Utsname
	err := syscall.Uname(&u)
	if err != nil {
		return "", err
	}
	return t.int65toString(&u.Machine), nil
}

// UpdateImageProperties - copy ImageFields to properties
func (t *ImageMetadata) UpdateImageProperties(im *srv.ImageFields, properties map[any]any) {
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
