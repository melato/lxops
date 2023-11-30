package doc

import (
	"fmt"
	"io/fs"
	"reflect"
	"regexp"
	"strings"

	"melato.org/lxops/internal/util"

	"melato.org/lxops/yaml"
)

type Doc struct {
	FS              fs.FS `name:"-"`
	All             bool  `usage:"include undocumented fields"`
	Hidden          bool  `usage:"show only undocumented fields"`
	typeDescriptors map[string]*TypeDescriptor
}

func (t *Doc) Configured() error {
	t.typeDescriptors = make(map[string]*TypeDescriptor)
	return nil
}

func (t *Doc) readTypeDescriptor(name string) (*TypeDescriptor, error) {
	file := name + ".yaml"
	data, err := fs.ReadFile(t.FS, file)
	if err != nil {
		return nil, err
	}
	var descriptor TypeDescriptor
	err = yaml.Unmarshal(data, &descriptor)
	if err != nil {
		return nil, err
	}
	return &descriptor, nil
}

func (t *Doc) getTypeDescriptor(name string) (*TypeDescriptor, error) {
	desc, ok := t.typeDescriptors[name]
	if !ok {
		var err error
		desc, err = t.readTypeDescriptor(name)
		if err != nil {
			return nil, err
		}
		t.typeDescriptors[name] = desc
	}
	return desc, nil
}

var rePkg = regexp.MustCompile(`[^\]\.]+\.`)

func typeName(typ reflect.Type) string {
	return rePkg.ReplaceAllString(typ.String(), "")
}

func (t *Doc) printFields(typ reflect.Type, name string, indent int, printMaps bool) error {
	descriptor, err := t.getTypeDescriptor(name)
	if err != nil {
		return err
	}
	if typ.Kind() != reflect.Struct {
		return nil
	}
	n := typ.NumField()
	for i := 0; i < n; i++ {
		f := typ.Field(i)
		if !f.IsExported() {
			continue
		}
		if f.Type.Kind() == reflect.Struct {
			err := t.printFields(f.Type, name, indent, printMaps)
			if err != nil {
				return err
			}
		}
		parts := strings.SplitN(f.Tag.Get("yaml"), ",", 2)
		name := parts[0]

		if name == "-" || name == "" {
			continue
		}

		desc, ok := descriptor.Fields[name]
		if !ok && (!t.All && !t.Hidden) {
			continue
		}
		if t.Hidden && ok {
			continue
		}
		fmt.Printf("%*s%s: (%v)\n", indent*2, "", name, typeName(f.Type))
		lines := util.SplitLines(desc)
		for _, line := range lines {
			fmt.Printf("%*s%s\n", (indent+1)*2, "", line)
		}
		fmt.Println()
		if f.Type.Kind() == reflect.Map && printMaps {
			elem := f.Type.Elem()
			if elem.Kind() == reflect.Pointer {
				elem = elem.Elem()
			}
			if elem.Kind() == reflect.Struct {
				err := t.printFields(elem, elem.Name(), indent+1, printMaps)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (t *Doc) printType(typ reflect.Type, name string) error {
	descriptor, err := t.getTypeDescriptor(name)
	if err != nil {
		return err
	}
	fmt.Printf("%s:\n", name)
	lines := util.SplitLines(descriptor.Description)
	for _, line := range lines {
		fmt.Printf("%*s%s\n", 2, "", line)
	}
	fmt.Println()
	return t.printFields(typ, name, 0, false)
}

func (t *Doc) PrintType(pointer any, name string) error {
	typ := reflect.TypeOf(pointer).Elem()
	return t.printType(typ, name)
}

func (t *Doc) PrintTypeFunc(pointer any, name string) func() error {
	return func() error {
		return t.PrintType(pointer, name)
	}
}
