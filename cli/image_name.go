package cli

import (
	"fmt"
	"regexp"
	"time"
)

type PrefixNameTime struct {
	Prefix string
	Name   string
	Time   time.Time
}

func (t *PrefixNameTime) Parse(name string) error {
	re := regexp.MustCompile("^([^-]+)-([^-]+)-([0-9]{8}-[0-9]{4})$")
	parts := re.FindStringSubmatch(name)
	if len(parts) != 4 {
		return fmt.Errorf("cannot parse: %s parts=%d", name, len(parts))
	}
	t.Prefix = parts[1]
	t.Name = parts[2]
	tm, err := time.Parse("20060102-1504", parts[3])
	if err != nil {
		return err
	}
	t.Time = tm
	return nil
}

// ParsePrefixNameDateTime - parse a name of the form {variant}-{os}-{date}-{time}
// example: a-nginx-20240203-0834 -> os=nginx variant=a date/time=20240203-0834
func (t *ImageMetadataOptions) ParsePrefixNameDatetimeVariantOS(name string) error {
	var p PrefixNameTime
	err := p.Parse(name)
	if err != nil {
		return err
	}

	t.OS = p.Name
	t.Variant = p.Prefix
	t.Date = p.Time
	return nil
}

// ParsePrefixNameDateTime - parse a name of the form {prefix}-{name}-{date}-{time}
// example: a-nginx-20240203-0834.
// If format=vod, set variant=prefix, os=name
// If format=ovd, set os=prefix, variant=name
func (t *ImageMetadataOptions) ParsePrefixNameTime(name string, format string) error {
	var p PrefixNameTime
	err := p.Parse(name)
	if err != nil {
		return err
	}
	switch format {
	case "vod":
		t.Variant = p.Prefix
		t.OS = p.Name
		t.Date = p.Time
		return nil
	case "ovd":
		t.OS = p.Prefix
		t.Variant = p.Name
		t.Date = p.Time
		return nil
	default:
		fmt.Printf("valid parse values are:\n")
		fmt.Printf("vod: {variant}-{os}-{datetime}\n")
		fmt.Printf("ovd: {os}-{variant}-{datetime}\n")
		return fmt.Errorf("unsupported name format: %s", format)
	}
}
