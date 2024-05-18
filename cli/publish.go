package cli

import (
	"fmt"
	"reflect"
	"strings"

	"melato.org/lxops/srv"
	"melato.org/lxops/yaml"
)

// PublishOps - publish instances
type PublishOps struct {
	InstanceOps
	Fields srv.ImageFields
	Alias  string `name:"alias" usage:"image alias`
	DryRun bool   `name:"dry-run" usage:"show image properties without publishing`
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

func (t *PublishOps) parseInstanceSnapshot(instanceSnapshot string) (instance string, snapshot string, err error) {
	parts := strings.SplitN(instanceSnapshot, "/", 3)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("%s: expected <instance>/<snapshot name>", instanceSnapshot)
	}
	return parts[0], parts[1], nil
}

func (t *PublishOps) PublishInstance1(instanceSnapshot string) error {
	instance, snapshot, err := t.parseInstanceSnapshot(instanceSnapshot)
	if err != nil {
		return err
	}
	return t.PublishInstance(instance, snapshot)
}

func (t *PublishOps) PublishInstance(instance, snapshot string) error {
	im, err := t.server.GetInstanceImageFields(instance)
	if err != nil {
		return err
	}
	t.mergeStructs(im, &t.Fields)
	if t.Fields.Name == "" {
		im.Name = fmt.Sprintf("%s-%s-%s-%s-%s", im.OS, im.Release, im.Architecture, im.Variant, im.Serial)
	}
	if t.Fields.Description == "" {
		im.Description = fmt.Sprintf("%s %s %s (%s)", im.OS, im.Release, im.Architecture, im.Serial)
	}
	alias := t.Alias
	if alias == "" {
		alias = instance
	}
	yaml.Print(im)
	if !t.DryRun {
		return t.server.PublishInstanceWithFields(instance, snapshot, alias, *im)
	}
	return nil
}

func (t *InstanceOps) Info(instance string) error {
	i, err := t.server.GetInstance(instance)
	if err != nil {
		return err
	}
	return yaml.Print(i)
}
