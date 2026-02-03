package lxops

import (
	"fmt"
	"io/fs"
	"os"

	"melato.org/lxops/cfg"
	"melato.org/lxops/help"
	"melato.org/lxops/internal/templatefs"
)

func helpDataModel(serverType string) any {
	return map[string]string{
		"ServerType":    serverType,
		"ConfigVersion": cfg.Comment,
		"lxops_doc_url": "https://github.com/melato/lxops/blob/main/md",
	}
}

func helpFS(serverType string) fs.FS {
	var fsys fs.FS = help.FS
	helpDir, ok := os.LookupEnv("LXOPS_HELP")
	if ok {
		fsys = os.DirFS(helpDir)
	}
	return templatefs.NewTemplateFS(fsys, helpDataModel(serverType))
}

func getUsage(serverType string) []byte {
	helpFS := helpFS(serverType)
	usageData, err := fs.ReadFile(helpFS, "commands.yaml.tpl")
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	return usageData
}
