// package types is used for documentation, by the help command.
// It contains descriptions for various types and their fields.
package help

import (
	"embed"
)

//go:embed *.yaml
var FS embed.FS
