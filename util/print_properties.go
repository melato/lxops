package util

import (
	"os"
	"sort"

	"melato.org/table3"
)

func PrintProperties(properties map[string]string) {
	keys := make([]string, 0, len(properties))
	for key, _ := range properties {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	var key, value string
	writer := &table.FixedWriter{Writer: os.Stdout}
	writer.Columns(
		table.NewColumn("KEY", func() interface{} { return key }),
		table.NewColumn("VALUE", func() interface{} { return value }),
	)

	for _, key = range keys {
		value = properties[key]
		writer.WriteRow()
	}
	writer.End()
}
