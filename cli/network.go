package cli

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"

	"melato.org/lxops/srv"
)

const DefaultProject = "default"

func QualifiedContainerName(project string, container string) string {
	if project == DefaultProject {
		return container
	}
	return project + "_" + container
}

type NetworkManager struct {
	Client srv.Client
}

func (t *NetworkManager) ParseAddress(addr string) string {
	i := strings.Index(addr, " ")
	if i > 0 {
		return addr[0:i]
	}
	return ""
}
func (t *NetworkManager) GetAddresses(family string) ([]*srv.HostAddress, error) {
	projects, err := t.Client.Projects()
	if err != nil {
		return nil, err
	}
	var addresses []*srv.HostAddress
	for _, project := range projects {
		server, err := t.Client.ProjectInstanceServer(project)

		if err != nil {
			return nil, err
		}
		paddresses, err := server.GetInstanceAddresses(family)
		if err != nil {
			return nil, err
		}
		addresses = append(addresses, paddresses...)
	}
	return addresses, nil
}

func (t *NetworkManager) printAddresses(addresses []*srv.HostAddress, headers bool, writer io.Writer) error {
	var csv = csv.NewWriter(writer)
	if headers {
		csv.Write([]string{"address", "name"})
	}
	for _, a := range addresses {
		csv.Write([]string{a.Address, a.Name})
	}
	csv.Flush()
	return csv.Error()
}

type NetworkOp struct {
	Client     srv.Client `name:"-"`
	OutputFile string     `name:"o" usage:"output file"`
	Format     string     `name:"format" usage:"include format: csv | yaml"`
	Headers    bool       `name:"headers" usage:"include headers"`
	Family     string     `name:"family" usage:"network family: inet | inet6"`
}

func (t *NetworkOp) Init() error {
	t.Family = "inet"
	t.Format = "csv"
	return nil
}

func (t *NetworkOp) ExportAddresses() error {
	net := &NetworkManager{Client: t.Client}
	addresses, err := net.GetAddresses(t.Family)

	if err != nil {
		return err
	}

	var printer AddressPrinter
	switch t.Format {
	case "csv":
		printer = &CsvAddressPrinter{Headers: t.Headers}
	case "yaml":
		printer = &YamlAddressPrinter{}
	default:
		return fmt.Errorf("unrecognized format: %s", t.Format)
	}

	var out io.WriteCloser
	if t.OutputFile == "" {
		out = os.Stdout
	} else {
		out, err = os.Create(t.OutputFile)
		if err != nil {
			return err
		}
		defer out.Close()
	}
	return printer.Print(addresses, out)
}
