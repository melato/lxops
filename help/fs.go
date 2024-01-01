// package types is used for documentation, by the help command.
// It contains descriptions for various types and their fields.
package help

import (
	"embed"
)

//go:embed *.tpl *.yaml topics/*
var FS embed.FS
