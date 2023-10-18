// Package api specifies the interface between lxdops and the instance server (LXD, Incus).
package srv

type InstanceServer interface {
	// Profiles
	GetProfileNames() ([]string, error)
	CreateProfile(profile *Profile) error
	DeleteProfile(profile string) error
	ProfileExists(name string) (bool, error)
	ExportProfile(name string) ([]byte, error)
	ImportProfile(name string, data []byte) error

	// instances
	LaunchInstance(launch *Launch) error
	CopyInstance(cp *Copy) error
	RenameInstance(oldname, newname string) error
	NewConfigurer(instance string) (InstanceConfigurer, error)

	GetInstanceProfiles(instance string) ([]string, error)
	SetInstanceProfiles(instance string, profiles []string) error

	GetInstanceNetwork(name string) (Network, error)
	SetInstanceNetwork(name string, network Network) error

	StartInstance(name string) error
	StopInstance(name string) error
	DeleteInstance(name string, stop bool) error
	WaitForNetwork(instance string) error
	CreateInstanceSnapshot(instance string, snapshot string) error

	// images
	ExportImage(image string, path string) error

	// informational methods, may be removed.
	GetInstanceDevices(name string) (map[string]*Device, error)
	GetHwaddresses() ([]Hwaddr, error)
	GetInstanceImages() ([]InstanceImage, error)
	GetInstanceNames(onlyRunning bool) ([]string, error)
	// GetInstanceAddresses - family is "inet" or "inet6"
	GetInstanceAddresses(family string) ([]*HostAddress, error)
}

type Profile struct {
	Name        string
	Description string
	Devices     map[string]*Device
	Config      map[string]string
}

type Device struct {
	Path     string
	Source   string
	Readonly bool
}

type Network interface{}

type Launch struct {
	Name       string
	Image      string
	Project    string
	Profiles   []string
	LxcOptions []string
	Network    Network
}

type Copy struct {
	Name           string
	Project        string
	SourceInstance string
	SourceSnapshot string
	Network        Network
}

type Hwaddr struct {
	Instance string
	Hwaddr   string
}

type InstanceImage struct {
	Instance string
	Image    string
}

type HostAddress struct {
	Name    string
	Address string
}
