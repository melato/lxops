package cli

import (
	"fmt"
	"time"
	"unicode"

	"melato.org/lxops/yaml"
)

// PublishOps - publish instances
type PublishOps struct {
	InstanceOps
	Name  string `name:"name" usage:"simple image name (one word, lowercase)`
	Alias string `name:"alias" usage:"image alias`
}

// capFirst - capitalize first rune of string
func capFirst(s string) string {
	if len(s) == 0 {
		return ""
	}
	var s0 string
	for i, c := range s {
		if i == 0 {
			s0 = string(unicode.ToUpper(c))
		} else {
			return s0 + s[i:]
		}
	}
	return s0
}

func (t *PublishOps) PublishInstance(instance, snapshot string) error {
	im, err := t.server.GetInstanceImageFields(instance)
	if err != nil {
		return err
	}
	/*
	  image.description: Alpinelinux 3.19 x86_64 (20240129_1300)
	  image.name: alpinelinux-3.19-x86_64-default-20240129_1300
	  image.os: alpinelinux
	  image.release: "3.19"
	  image.serial: "20240129_1300"
	  image.variant: default
	*/
	serial := time.Now().UTC().Format("20060102_1504")
	im.Serial = serial
	name := t.Name
	if name == "" {
		name = instance
	}
	im.Variant = name
	alias := t.Alias
	if alias == "" {
		alias = instance
		//alias = name + "-" + strings.ReplaceAll(serial, "_", "-")
	}
	im.Name = fmt.Sprintf("%s-%s-%s-%s-%s", im.OS, im.Release, im.Architecture, im.Variant, im.Serial)
	im.Description = fmt.Sprintf("%s %s %s (%s)", capFirst(name), im.Release, im.Architecture, im.Serial)
	return t.server.PublishInstanceWithFields(instance, snapshot, alias, *im)
}

func (t *InstanceOps) Info(instance string) error {
	i, err := t.server.GetInstance(instance)
	if err != nil {
		return err
	}
	return yaml.Print(i)
}
