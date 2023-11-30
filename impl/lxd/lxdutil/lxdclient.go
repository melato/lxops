package lxdutil

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	lxd "github.com/canonical/lxd/client"
	"github.com/canonical/lxd/lxc/config"
	"melato.org/lxops/srv"
	"melato.org/lxops/yaml"
)

type LxdClient struct {
	Socket string
	Http   bool `usage:"connect to LXD using http"`
	Unix   bool `usage:"connect to LXD using unix socket"`
	//Project        string `name:"project" usage:"the LXD project to use.  Overrides Config.Project"`
	rootServer    lxd.InstanceServer
	projectServer lxd.InstanceServer
	LxcConfig
}

func (t *LxdClient) Init() error {
	sockets := []string{
		"/var/snap/lxd/common/lxd/unix.socket",
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

// connectUnix - Connect to LXD over the Unix socket
func (t *LxdClient) connectUnix() (lxd.InstanceServer, error) {
	server, err := lxd.ConnectLXDUnix(t.Socket, nil)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("%s: %s", t.Socket, err.Error()))
	}
	return server, nil
}

func (t *LxdClient) configFile(name string) (string, error) {
	dir, err := ConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, name), nil
}

func (t *LxdClient) readConfigFile(name string) ([]byte, error) {
	file, err := t.configFile(name)
	if err != nil {
		return nil, err
	}
	return os.ReadFile(file)
}

func (t *LxdClient) connectHttp() (lxd.InstanceServer, error) {
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
	args := &lxd.ConnectionArgs{
		AuthType:      remote.AuthType,
		TLSServerCert: string(serverCrt),
		TLSClientCert: string(crt),
		TLSClientKey:  string(key)}
	server, err := lxd.ConnectLXD(remote.Addr, args)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("%s: %s", remote.Addr, err.Error()))
	}
	return server, nil
}

func (t *LxdClient) RootServer() (lxd.InstanceServer, error) {
	if t.rootServer == nil {
		var server lxd.InstanceServer
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

func (t *LxdClient) Projects() ([]string, error) {
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

func (t *LxdClient) ProjectServer(project string) (lxd.InstanceServer, error) {
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

func (t *LxdClient) CurrentServer() (lxd.InstanceServer, error) {
	return t.ProjectServer("")
}

func (t *LxdClient) CurrentInstanceServer() (srv.InstanceServer, error) {
	return t.ProjectInstanceServer("")
}

func (t *LxdClient) ProjectInstanceServer(project string) (srv.InstanceServer, error) {
	server, err := t.ProjectServer(project)
	if err != nil {
		return nil, err
	}
	return &InstanceServer{Server: server}, nil
}

func (t *LxdClient) ServerType() string {
	return "LXD"
}
