package cli

import (
	"fmt"
	"regexp"
	"time"
)

// ParsePrefixNameDateTime - parse a name of the form {variant}-{os}-{date}-{time}
// example: a-nginx-20240203-0834 -> os=nginx variant=a date/time=20240203-0834
func (t *ImageMetadataOptions) ParsePrefixNameDateTime(name string) error {
	re := regexp.MustCompile("^([^-]+)-([^-]+)-([0-9]{8}-[0-9]{4})$")
	parts := re.FindStringSubmatch(name)
	if len(parts) != 4 {
		return fmt.Errorf("cannot parse: %s parts=%d", name, len(parts))
	}
	t.Variant = parts[1]
	t.OS = parts[2]
	tm, err := time.Parse("20060102-1504", parts[3])
	if err != nil {
		return nil
	}
	t.Date = tm
	return nil
}
