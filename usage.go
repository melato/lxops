package lxops

import (
	"bytes"
	_ "embed"
	"fmt"
	"text/template"

	"melato.org/lxops/cfg"
)

//go:embed commands.yaml.tpl
var usageTemplate string

func helpDataModel(serverType string) any {
	return map[string]string{
		"ServerType":    serverType,
		"ConfigVersion": cfg.Comment,
		"lxops_doc_url": "https://github.com/melato/lxops/blob/main/md",
	}
}

func getUsage(serverType string) []byte {
	tpl, err := template.New("").Parse(usageTemplate)
	var buf bytes.Buffer
	if err == nil {
		model := helpDataModel(serverType)
		err = tpl.Execute(&buf, model)
	}
	if err != nil {
		fmt.Printf("%v\n", err)
		return nil
	}
	return buf.Bytes()
}
