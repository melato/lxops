package cli

import (
	"encoding/csv"
	"io"

	"melato.org/lxops/srv"
	"melato.org/lxops/yaml"
)

type AddressPrinter interface {
	Print(addresses []*srv.HostAddress, writer io.Writer) error
}

type CsvAddressPrinter struct {
	Headers bool
}

func (t *CsvAddressPrinter) Print(addresses []*srv.HostAddress, writer io.Writer) error {
	var csv = csv.NewWriter(writer)
	if t.Headers {
		csv.Write([]string{"address", "name"})
	}
	for _, a := range addresses {
		csv.Write([]string{a.Address, a.Name})
	}
	csv.Flush()
	return csv.Error()
}

type YamlAddressPrinter struct {
}

func (t *YamlAddressPrinter) Print(addresses []*srv.HostAddress, writer io.Writer) error {
	m := make(map[string]string)
	for _, a := range addresses {
		m[a.Name] = a.Address
	}
	data, err := yaml.Marshal(m)
	if err == nil {
		_, err = writer.Write(data)
	}
	return err
}
