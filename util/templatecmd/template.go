package templatecmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"text/template"

	"melato.org/lxops/util/templatecmd/property"
	"melato.org/lxops/yaml"
)

type TemplateOp struct {
	TemplateFile string   `name:"t" usage:"template file"`
	OutputFile   string   `name:"o" usage:"output file"`
	FileMode     string   `name:"mode" usage:"file mode"`
	KeyFiles     []string `name:"F" usage:"key=file - set the value of <key> to the content of <file>"`
	property.Options
	PrintKeys   bool             `name:"print-keys" usage:"print keys without applying template"`
	PrintValues bool             `name:"print-values" usage:"print key/values without applying template"`
	Funcs       template.FuncMap `name:"-"`
}

func (t *TemplateOp) Init() error {
	t.FileMode = "0664"
	return nil
}

func (t *TemplateOp) Configured() error {
	if t.TemplateFile == "" {
		return errors.New("missing template file")
	}
	return nil
}

func parsePairs(args []string) (map[string]string, error) {
	result := make(map[string]string)
	for _, arg := range args {
		eq := strings.Index(arg, "=")
		if eq < 0 {
			return nil, errors.New("invalid key-pair: " + arg)
		}
		key := arg[:eq]
		value := arg[eq+1:]
		result[key] = value
	}
	return result, nil
}

func (t *TemplateOp) Values() (property.Properties, error) {
	properties := make(property.Properties)
	err := t.Options.Apply(properties)
	if err != nil {
		return nil, err
	}
	files, err := parsePairs(t.KeyFiles)
	if err != nil {
		return nil, err
	}
	for key, file := range files {
		value, err := os.ReadFile(file)
		if err != nil {
			return nil, err
		}
		properties[key] = string(value)
	}
	return properties, nil
}

func (t *TemplateOp) Run() error {
	values, err := t.Values()
	if err != nil {
		return err
	}
	if t.PrintKeys {
		keys := make([]string, 0, len(values))
		var otherKeys []interface{}
		for key, _ := range values {
			s, isString := key.(string)
			if isString {
				keys = append(keys, s)
			} else {
				otherKeys = append(otherKeys, key)
			}
		}
		sort.Strings(keys)
		if t.PrintKeys {
			for _, key := range keys {
				fmt.Printf("%s\n", key)
			}
		}
		for _, key := range otherKeys {
			fmt.Printf("%v\n", key)
		}
	}
	if t.PrintValues {
		yaml.Print(values)
	}

	if t.PrintKeys || t.PrintValues {
		return nil
	}

	tpl0 := template.New("x")
	tpl0.Funcs(t.Funcs)
	_, err = tpl0.ParseFiles(t.TemplateFile)
	if err != nil {
		return err
	}
	tpl := tpl0.Lookup(filepath.Base(t.TemplateFile))
	if tpl == nil {
		fmt.Printf("%s\n", tpl0.DefinedTemplates())
		return fmt.Errorf("could not lookup template")
	}

	w := os.Stdout
	var tmpName string
	if t.OutputFile != "" {
		mode, err := strconv.ParseInt(t.FileMode, 8, 32)
		if err != nil {
			return fmt.Errorf("invalid mode: %s", t.FileMode)
		}
		dir := filepath.Dir(t.OutputFile)
		w, err = os.CreateTemp(dir, "tpl")
		if err != nil {
			return nil
		}
		tmpName = w.Name()
		defer os.Remove(tmpName)
		err = os.Chmod(tmpName, os.FileMode(mode))
		if err != nil {
			return nil
		}
		w, err = os.OpenFile(tmpName, os.O_RDWR, os.FileMode(mode))
		if err != nil {
			return err
		}
		defer w.Close()
	}
	err = tpl.Execute(w, values)
	if t.OutputFile != "" {
		w.Close()
		if err != nil {
			return err
		}
		err = os.Rename(tmpName, t.OutputFile)
	}
	return err
}
