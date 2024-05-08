package lxops_incus

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	incus "github.com/lxc/incus/v6/client"
	config "github.com/lxc/incus/v6/shared/cliconfig"
	"melato.org/lxops/srv"
	"melato.org/lxops/yaml"
)

type Client struct {
	Socket string
	Http   bool `usage:"connect to Incus using http"`
	Unix   bool `usage:"connect to Incus using unix socket"`
	//Project        string `name:"project" usage:"the Incus project to use.  Overrides Config.Project"`
	rootServer    incus.InstanceServer
	projectServer incus.InstanceServer
	Config
}

func (t *Client) Init() error {
	sockets := []string{
		"/var/lib/incus/unix.socket",
	}
	for _, socket := range sockets {
		_, err := os.Stat(socket)
		if err == nil {
			t.Socket = socket
			break
		}
	}
	return nil
}

// connectUnix - Connect to Incus over the Unix socket
func (t *Client) connectUnix() (incus.InstanceServer, error) {
	server, err := incus.ConnectIncusUnix(t.Socket, nil)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("%s: %s", t.Socket, err.Error()))
	}
	return server, nil
}

func (t *Client) configFile(name string) (string, error) {
	dir, err := ConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, name), nil
}

func (t *Client) readConfigFile(name string) ([]byte, error) {
	file, err := t.configFile(name)
	if err != nil {
		return nil, err
	}
	return os.ReadFile(file)
}

func (t *Client) connectHttp() (incus.InstanceServer, error) {
	var cfg config.Config
	cfgPath, err := t.configFile("config.yml")
	if err != nil {
		return nil, err
	}
	err = yaml.ReadFile(cfgPath, &cfg)
	if err != nil {
		return nil, err
	}
	if cfg.DefaultRemote == "" {
		return nil, fmt.Errorf("missing default remote")
	}
	remote, found := cfg.Remotes[cfg.DefaultRemote]
	if !found {
		return nil, fmt.Errorf("missing remote: %s", cfg.DefaultRemote)
	}
	serverCrt, err := t.readConfigFile(fmt.Sprintf("servercerts/%s.crt", cfg.DefaultRemote))
	if err != nil {
		return nil, err
	}
	crt, err := t.readConfigFile("client.crt")
	if err != nil {
		return nil, err
	}
	key, err := t.readConfigFile("client.key")
	if err != nil {
		return nil, err
	}
	args := &incus.ConnectionArgs{
		AuthType:      remote.AuthType,
		TLSServerCert: string(serverCrt),
		TLSClientCert: string(crt),
		TLSClientKey:  string(key)}
	server, err := incus.ConnectIncus(remote.Addr, args)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("%s: %s", remote.Addr, err.Error()))
	}
	return server, nil
}

func (t *Client) RootServer() (incus.InstanceServer, error) {
	if t.rootServer == nil {
		var server incus.InstanceServer
		var err error
		if t.Http {
			server, err = t.connectHttp()
		} else if t.Unix {
			server, err = t.connectUnix()
		} else {
			server, err = t.connectHttp()
			if err != nil {
				server, err = t.connectUnix()
			}
		}
		if err != nil {
			return nil, err
		}
		t.rootServer = server
	}
	return t.rootServer, nil
}

func (t *Client) Projects() ([]string, error) {
	server, err := t.RootServer()
	if err != nil {
		return nil, err
	}
	projects, err := server.GetProjects()
	if err != nil {
		return nil, err
	}
	names := make([]string, len(projects))
	for i, project := range projects {
		names[i] = project.Name
	}
	return names, nil
}

func (t *Client) ProjectServer(project string) (incus.InstanceServer, error) {
	var err error
	if project == "" {
		project = t.CurrentProject()
	}
	server, err := t.RootServer()
	if err != nil {
		return nil, err
	}
	if project == "default" {
		return server, nil
	}
	return server.UseProject(project), nil
}

func (t *Client) CurrentServer() (incus.InstanceServer, error) {
	return t.ProjectServer("")
}

func (t *Client) CurrentInstanceServer() (srv.InstanceServer, error) {
	return t.ProjectInstanceServer("")
}

func (t *Client) ProjectInstanceServer(project string) (srv.InstanceServer, error) {
	server, err := t.ProjectServer(project)
	if err != nil {
		return nil, err
	}
	return &InstanceServer{Server: server}, nil
}

func (t *Client) ServerType() string {
	return "Incus"
}
