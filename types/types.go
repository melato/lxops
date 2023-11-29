// package types is used for documentation, by the help command.
// It contains descriptions for various types and their fields.
package types

import (
	"embed"
)

//go:embed *.yaml
var Types embed.FS
