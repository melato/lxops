package cli

import (
	"fmt"
	"time"

	"melato.org/lxops/srv"
	"melato.org/lxops/yaml"
)

// ImageMetadata - represents image metadata.yaml, in generic form
type ImageMetadata struct {
	v          map[any]any
	properties map[any]any
}

func ReadImageMetadata(file string) (*ImageMetadata, error) {
	var v map[any]any
	var properties map[any]any
	err := yaml.ReadFile(file, &v)
	if err != nil {
		return nil, err
	}
	vProperties := v["properties"]
	if vProperties == nil {
		return nil, fmt.Errorf("missing properties")
	}
	var ok bool
	properties, ok = vProperties.(map[any]any)
	if !ok {
		return nil, fmt.Errorf("properties: %T. expected: map[any]any", vProperties)
	}
	var t ImageMetadata
	t.v = v
	t.properties = properties
	return &t, nil
}

func (t *ImageMetadata) WriteFile(file string) error {
	return yaml.WriteFile(t.v, file)
}

func (t *ImageMetadata) Print() error {
	return yaml.Print(t.v)
}

func NewImageMetadata() (*ImageMetadata, error) {
	var t ImageMetadata
	t.v = make(map[any]any)
	t.properties = make(map[any]any)
	t.v["properties"] = t.properties
	var err error
	arch, err := getSystemArchitecture()
	if err != nil {
		return nil, err
	}
	t.properties["architecture"] = arch
	return &t, nil
}

func (t *ImageMetadata) getProperty(name string) string {
	v, ok := t.properties[name]
	if ok {
		s, ok := v.(string)
		if ok {
			return s
		}
	}
	return ""
}

func (t *ImageMetadata) GetFields() *srv.ImageFields {
	var f srv.ImageFields
	f.Architecture = t.getProperty("architecture")
	f.Description = t.getProperty("description")
	f.Name = t.getProperty("name")
	f.OS = t.getProperty("os")
	f.Release = t.getProperty("release")
	f.Serial = t.getProperty("serial")
	f.Variant = t.getProperty("variant")
	return &f
}

// SetFields - copy ImageFields to properties
func (t *ImageMetadata) SetFields(im *srv.ImageFields) {
	update := func(name, value string) {
		if value != "" {
			t.properties[name] = value
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

// SetFields - set creation_date, expiry_date
func (t *ImageMetadata) SetDates(date time.Time, expiryDays int) {
	timestamp := date.Unix()
	t.v["creation_date"] = timestamp
	t.v["expiry_date"] = timestamp + int64(expiryDays)*24*3600
}
