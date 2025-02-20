package cli

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"

	"melato.org/lxops/srv"
	"melato.org/lxops/yaml"
)

// PublishOps - publish instances
type PublishOps struct {
	InstanceOps
	Fields srv.ImageFields
	Alias  string `name:"alias" usage:"image alias"`
}

func (t *PublishOps) Init() error {
	t.Fields.Architecture = runtime.GOARCH
	return nil
}

func (_ *PublishOps) mergeStructs(target, source any) {
	t := reflect.TypeOf(source).Elem()
	vIn := reflect.ValueOf(source).Elem()
	vOut := reflect.ValueOf(target).Elem()
	n := t.NumField()
	for i := 0; i < n; i++ {
		f := t.Field(i)
		in := vIn.Field(i)
		if !in.IsZero() {
			//fmt.Printf("set %s=%v\n", f.Name, in)
			out := vOut.FieldByName(f.Name)
			out.Set(in)
		}
	}
}

func parseInstanceSnapshot(instanceSnapshot string) (instance string, snapshot string, err error) {
	parts := strings.SplitN(instanceSnapshot, "/", 3)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("%s: expected <instance>/<snapshot name>", instanceSnapshot)
	}
	return parts[0], parts[1], nil
}

func (t *PublishOps) PublishInstance(instanceSnapshot string) error {
	instance, snapshot, err := parseInstanceSnapshot(instanceSnapshot)
	if err != nil {
		return err
	}
	return t.PublishInstance2(instance, snapshot)
}

func (t *PublishOps) PublishInstance2(instance, snapshot string) error {
	alias := t.Alias
	if alias == "" {
		alias = instance
	}
	return t.server.PublishInstanceWithFields(instance, snapshot, alias, t.Fields)
}

func (t *InstanceOps) Info(instance string) error {
	i, err := t.server.GetInstance(instance)
	if err != nil {
		return err
	}
	return yaml.Print(i)
}
