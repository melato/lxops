package cli

import (
	"fmt"
	"syscall"
	"time"

	"melato.org/lxops/srv"
)

// ImageMetadataOptions - independent image properties
type ImageMetadataOptions struct {
	OS          string    `name:"os" usage:"override os property"`
	Variant     string    `name:"variant" usage:"override variant property"`
	Release     string    `name:"release" usage:"override release property, use '.' for timestamp"`
	Serial      string    `name:"serial" usage:"override serial property, use '.' for timestamp"`
	Name        string    `name:"name" usage:"override name property"`
	Description string    `name:"description" usage:"override description property"`
	Date        time.Time `name:"-"`
	DateString  string    `name:"date" usage:"override creation_date as YYYYMMDD-HHMM. use '.' for current time"`
	ExpiryDays  int       `name:"expiry-days" usage:"expiry_date as number of days from creation_date"`
}

func (t *ImageMetadataOptions) Init() error {
	t.ExpiryDays = 30
	return nil
}

func (t *ImageMetadataOptions) Configured() error {
	now := time.Now()
	if t.DateString != "" {
		if t.DateString == "." {
			t.Date = now
		} else {
			tm, err := time.Parse("20060102-1504", t.DateString)
			if err != nil {
				return err
			}
			t.Date = tm
		}
	}
	timestamp := now.UTC().Format("200601021504")
	if t.Release == "." {
		t.Release = timestamp
	}
	if t.Serial == "." {
		t.Serial = timestamp
	}
	return nil
}

// SetImageDescription - generate name, description from other fields
func (t *ImageMetadataOptions) SetImageDescription(im *srv.ImageFields, name string) {
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
func (t *ImageMetadataOptions) UpdateImageProperties(im *srv.ImageFields, properties map[any]any) {
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

func (t *ImageMetadataOptions) Override(im *srv.ImageFields) {
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
	if t.Serial != "" {
		im.Serial = t.Serial
	}
	if t.Name != "" {
		im.Name = t.Name
	}
	if t.Description != "" {
		im.Description = t.Description
	}
}

func (t *ImageMetadataOptions) HasChanges() bool {
	return t.Release != "" ||
		t.OS != "" ||
		t.Variant != "" ||
		t.Serial != "" ||
		t.Release != "" ||
		t.Name != "" ||
		t.Description != "" ||
		!t.Date.IsZero()
}

func (t *ImageMetadataOptions) Apply(m *ImageMetadata) {
	f := m.GetFields()
	t.Override(f)
	t.SetImageDescription(f, t.Variant)
	m.SetFields(f)
	m.UpdateDates(t.Date, t.ExpiryDays)
}
