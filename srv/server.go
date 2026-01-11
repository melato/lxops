// Package api specifies the interface between lxops and the instance server (Incus, LXD).
package srv

type InstanceServer interface {
	// Profiles
	GetProfileNames() ([]string, error)
	GetProfileDevices(profile string) (map[string]*Device, error)
	CreateProfile(profile *Profile) error
	DeleteProfile(profile string) error
	ProfileExists(name string) (bool, error)
	ExportProfile(name string) ([]byte, error)
	ImportProfile(name string, data []byte) error

	GetStoragePool(pool string) (*StoragePool, error)

	// instances
	CreateInstance(launch *Create) error
	LaunchInstance(launch *Launch) error
	RebuildInstance(image, instance string) error
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
	PublishInstance(instance, snapshot, alias string) error
	PublishInstance2(instance, snapshot, alias string, fields ImageFields, options PublishOptions) error
	ExportImage(image string, path string) error
	ImportImage(image string, path string) error

	// informational methods, may be removed.
	GetInstanceDevices(name string) (map[string]*Device, error)
	GetHwaddresses() ([]Hwaddr, error)
	GetInstanceImages() ([]InstanceImage, error)
	GetInstanceImageFields(instance string) (*ImageFields, error)
	GetInstanceNames(onlyRunning bool) ([]string, error)
	// GetInstanceAddresses - family is "inet" or "inet6"
	GetInstanceAddresses(family string) ([]*HostAddress, error)

	GetInstance(name string) (any, error)
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
	Pool     string
}

type StoragePool struct {
	Name   string
	Driver string
	Source string
}

type Network interface{}

type Create struct {
	Name       string
	Image      string
	Project    string
	Profiles   []string
	LxcOptions []string
}

type Launch struct {
	Create
	Network Network
}

type Copy struct {
	Name           string
	Project        string
	SourceInstance string
	SourceSnapshot string
	Profiles       []string
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

type ImageFields struct {
	Architecture string
	OS           string
	Variant      string
	Release      string
	Serial       string
	Description  string
	Name         string
}

type PublishOptions struct {
	Compression string
	Format      string
}
