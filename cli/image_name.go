package cli

import (
	"fmt"
	"regexp"
	"time"
)

type PrefixNameTime struct {
	Prefix    string
	Name      string
	Timestamp string
}

func (t *PrefixNameTime) Parse(name string) error {
	re := regexp.MustCompile("^([^-]+)-([^-]+)-([0-9]{8}-[0-9]{4})$")
	parts := re.FindStringSubmatch(name)
	if len(parts) != 4 {
		return fmt.Errorf("cannot parse: %s parts=%d", name, len(parts))
	}
	t.Prefix = parts[1]
	t.Name = parts[2]
	t.Timestamp = parts[3]
	return nil
}

func (t *PrefixNameTime) ParseTime() (time.Time, error) {
	return time.Parse("20060102-1504", t.Timestamp)
}

// ParsePrefixNameDateTime - parse a name of the form {prefix}-{name}-{date}-{time}
// example: a-nginx-20240203-0834.
func (t *ImageMetadataOptions) ParsePrefixNameTime(name string, format string) error {
	var p PrefixNameTime
	err := p.Parse(name)
	if err != nil {
		return err
	}
	switch format {
	case "unique":
		t.Variant = t.Name
		t.Serial = p.Timestamp
		t.Release = p.Timestamp
		return nil
	case "reuse":
		t.Variant = t.Name
		t.Serial = p.Timestamp
		return nil
	case "ovr":
		t.OS = p.Prefix
		t.Variant = p.Name
		t.Release = p.Timestamp
		t.Serial = p.Timestamp
		return nil
	case "vod":
		t.Variant = p.Prefix
		t.OS = p.Name
		t.Serial = p.Timestamp
		t.Date, err = p.ParseTime()
		return err
	case "ovd":
		t.OS = p.Prefix
		t.Variant = p.Name
		t.Serial = p.Timestamp
		t.Date, err = p.ParseTime()
		return err
	default:
		return fmt.Errorf("unsupported name format: %s", format)
	}
}
