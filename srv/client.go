package srv

type Client interface {
	CurrentInstanceServer() (InstanceServer, error)
	ProjectInstanceServer(project string) (InstanceServer, error)
	Projects() ([]string, error)
	CurrentProject() string
}
