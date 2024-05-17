package cli

import (
	"fmt"
	"reflect"

	"melato.org/lxops/srv"
	"melato.org/lxops/yaml"
)

// PublishOps - publish instances
type PublishOps struct {
	InstanceOps
	Fields srv.ImageFields
	Alias  string `name:"alias" usage:"image alias`
}

func (_ *PublishOps) mergeStructs(target, source any) {
	t := reflect.TypeOf(source).Elem()
	vIn := reflect.ValueOf(source).Elem()
	vOut := reflect.ValueOf(target).Elem()
	n := t.NumField()
	fmt.Printf("%d fields\n", n)
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

func (t *PublishOps) PublishInstance(instance, snapshot string) error {
	im, err := t.server.GetInstanceImageFields(instance)
	if err != nil {
		return err
	}
	t.mergeStructs(im, &t.Fields)
	if im.Name == "" {
		im.Name = fmt.Sprintf("%s-%s-%s-%s-%s", im.OS, im.Release, im.Architecture, im.Variant, im.Serial)
	}
	if im.Description == "" {
		im.Description = fmt.Sprintf("%s %s %s (%s)", im.OS, im.Release, im.Architecture, im.Serial)
	}
	alias := t.Alias
	if alias == "" {
		alias = instance
	}
	return t.server.PublishInstanceWithFields(instance, snapshot, alias, *im)
}

func (t *InstanceOps) Info(instance string) error {
	i, err := t.server.GetInstance(instance)
	if err != nil {
		return err
	}
	return yaml.Print(i)
}
