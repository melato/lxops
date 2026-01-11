package lxops

import (
	"fmt"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"melato.org/lxops/cfg"
	"melato.org/lxops/cli"
	"melato.org/lxops/srv"
	"melato.org/lxops/yaml"
)

// PublishOps - publish instances
type PublishOps struct {
	cli.InstanceOps
	Fields     srv.ImageFields
	Options    srv.PublishOptions
	Alias      string `name:"alias" usage:"image alias"`
	ConfigFile string `name:"c" usage:"optional config file"`
	DryRun     bool   `name:"dry-run" usage:"print image properties but don't publish"`
	config     *cfg.Config
}

func (t *PublishOps) Init() error {
	t.Fields.Architecture = runtime.GOARCH
	return nil
}

func (t *PublishOps) parseVersion(name string) string {
	re := regexp.MustCompile(`[0-9][-0-9_\.]*[0-9]`)
	return re.FindString(name)
}

func (t *PublishOps) Configured() error {
	if t.ConfigFile != "" {
		var properties PropertyOptions
		err := properties.Init()
		if err == nil {
			properties.Configured()
		}
		reader := properties.NewConfigReader()
		t.config, err = reader.Read(t.ConfigFile)
		if err != nil {
			return err
		}
		if t.Fields.Description == "" {
			t.Fields.Description = t.config.Description
		}
		if t.Fields.OS == "" {
			t.Fields.OS = t.config.Ostype
		}
	}
	if !t.DryRun {
		return t.InstanceOps.Configured()
	}
	return nil
}

func parseInstanceSnapshot(instanceSnapshot string) (instance string, snapshot string, err error) {
	parts := strings.SplitN(instanceSnapshot, "/", 3)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("%s: expected <instance>/<snapshot name>", instanceSnapshot)
	}
	return parts[0], parts[1], nil
}

func join(args ...string) string {
	k := 0
	for _, arg := range args {
		if arg != "" {
			args[k] = arg
			k += 1
		}
	}
	return strings.Join(args[0:k], "-")
}

func (t *PublishOps) PublishInstance(instance string, args ...string) error {
	var snapshot string
	switch len(args) {
	case 1:
		snapshot = args[0]
	case 0:
		name, snap, hasSnapshot := strings.Cut(instance, "/")
		if hasSnapshot {
			instance = name
			snapshot = snap
		}
	default:
		return fmt.Errorf("too many arguments")
	}
	if t.config != nil {
		name := filepath.Base(t.ConfigFile)
		name = strings.TrimSuffix(name, filepath.Ext(name))
		if instance == "" {
			instance = name
		}
		if t.Fields.Variant == "" {
			t.Fields.Variant = name
		}
		if snapshot == "" {
			snapshot = t.config.Snapshot
		}
	}
	if instance == "" {
		return fmt.Errorf("need instance name or config file")
	}
	if snapshot == "" {
		return fmt.Errorf("missing snapshot")
	}

	if t.Fields.Variant == "" {
		t.Fields.Variant = instance
	}

	version := t.Fields.Serial
	if version == "" {
		version = t.Fields.Release
	}
	if version == "" {
		version = t.parseVersion(instance)
	}

	alias := t.Alias
	if alias == "" {
		if t.config != nil {
			alias = join(t.Fields.OS, t.Fields.Variant, version)
		} else {
			alias = instance
		}
	}
	if t.Fields.Serial == "" {
		t.Fields.Serial = version
	}
	if t.Fields.Release == "" {
		t.Fields.Release = version
	}

	if t.Fields.Name == "" {
		//t.Fields.Name = join(t.Fields.OS, instance, t.Fields.Release)
		t.Fields.Name = instance
	}

	fmt.Printf("instance: %s\n", instance)
	fmt.Printf("snapshot: %s\n", snapshot)
	fmt.Printf("alias: %s\n", alias)
	yaml.Print(&t.Options)
	yaml.Print(&t.Fields)
	if t.DryRun {
		return nil
	}
	return t.Server.PublishInstance2(instance, snapshot, alias, t.Fields, t.Options)
}
