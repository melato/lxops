package srv

type Client interface {
	CurrentInstanceServer() (InstanceServer, error)
	ProjectInstanceServer(project string) (InstanceServer, error)
	Projects() ([]string, error)
	CurrentProject() string
	// ServerType returns a string that identifies the server type.
	// This is used only for documentation.
	ServerType() string
}
