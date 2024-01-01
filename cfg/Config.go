package cfg

// HostPath is a file path on the host, which is either absolute or relative to a base directory
// When a config file includes another config file, the base directory is the directory of the including file
type HostPath string

// Config - Instance configuration
// The fields are documented better in the "lxops help" commands.
type Config struct {
	// File is the file of the top config file, for reference.  It is not read from yaml.
	File string `yaml:"-"`
	// ConfigInherit fields are merged with all included files, depth first
	ConfigInherit `yaml:",inline"`
	// ConfigTop fields are not merged with included files
	ConfigTop `yaml:",inline"`
}

type ConfigTop struct {
	// Description is provided for documentation
	Description string `yaml:"description,omitempty"`

	// Stop specifies that the container should be stopped at the end of the configuration
	Stop bool `yaml:"stop,omitempty"`

	// Snapshot specifies that that the container should be snapshoted with this name at the end of the configuration process.
	Snapshot string `yaml:"snapshot,omitempty"`
}

type ConfigInherit struct {
	// ostype - OS type for cloudconfig.  "alpine", "debian", etc.
	Ostype string `yaml:"ostype,omitempty"`

	// image - the image name (with optional remote), used when launching a container.
	Image Pattern `yaml:"image,omitempty"`

	// Project is the LXD or Incus project where the container is
	Project string `yaml:"project,omitempty"`

	// profile-config (deprecated)
	// ProfileConfig map[string]string `yaml:"profile-config,omitempty"`

	// Source specifies where to copy or clone the instance from
	Source `yaml:",inline"`

	// Extra options passed to lxc launch.
	LxcOptions []string `yaml:"lxc-options,omitempty,flow"`

	// profiles - the instance profiles
	Profiles []string `yaml:"profiles"`

	// ProfilePattern specifies how the instance profile should be named.
	// It defaults to "(instance).lxdops"
	Profile Pattern `yaml:"profile-pattern,omitempty"`

	// The owner (uid:gid) for new devices
	DeviceOwner Pattern `yaml:"device-owner,omitempty"`

	// Filesystems are zfs filesystems or plain directories that are created
	// when an instance is created.  Devices are created inside filesystems.
	Filesystems map[string]*Filesystem `yaml:"filesystems,omitempty"`

	// Devices are disk devices that are directories within the instance filesystems
	// They can also be standalone, without a filesystem.
	// They are created and attached to the container via the instance profile
	Devices map[string]*Device `yaml:"devices,omitempty"`
	// Profiles are attached to the container.  The instance profile should not be listed here.

	// Properties provide key-value pairs used for pattern substitution.
	Properties map[string]string `yaml:"properties,omitempty"`

	// Include is a list of other configs that are to be included.
	// Include paths are either absolute or relative to the path of the including config.
	Include []HostPath `yaml:"include,omitempty"`

	// cloud-config-files is a list of cloud-config files to run during instance configuration.
	CloudConfigFiles []HostPath `yaml:"cloud-config-files,omitempty"`
}

// Source specifies how to copy or clone the instance container, filesystem, and device directories.
// When DeviceTemplate is specified, the filesystems are copied with rsync.
// When DeviceOrigin is specified, the filesystems are cloned with zfs-clone
// The filesystems that are copied are determined by applying the source instance name to the filesystems of this config,
// or to the filesystems of a source config.
//
// When basing an instance on a template with few skeleton files, it is preferable to copy with a DeviceTemplate,
// so the container's disk devices are not tied to the template.
//
// Example:
// suppose test-a.yaml has:
//
//	origin: a/copy
//	filesystems: "default": "z/test/(instance)"
//	device-origin: a@copy
//	source-filesystems "default": "z/prod/(instance)"
//	devices: home, path=/home, filesystem=default
//
// This would do something like:
//
//	zfs clone z/prod/a@copy z/test/test-a
//	lxc copy --container-only a/copy test-a
//	lxc profile create test-a.lxdops
//	lxc profile device add test-a.lxdops home disk path=/home source=/z/test/test-a/home
//	lxc profile add test-a test-a.lxdops
type Source struct {
	// origin is the name of a container and a snapshot to clone from.
	// It has the form [<project>_]<container>[/<snapshot>]
	// It overrides SourceConfig
	Origin Pattern `yaml:"origin,omitempty"`

	// device-template is the name of an instance, whose devices are copied (using rsync)
	// to a new instance with launch.
	// The devices are copied from the filesystems specified in SourceConfig, or this config.
	DeviceTemplate Pattern `yaml:"device-template,omitempty"`

	// device-origin is the name an instance and a short snapshot name.
	// It has the form <instance>@<snapshot> where <instance> is an instance name,
	// and @<snapshot> is a the short snapshot name of the instance filesystems.
	// Each device zfs filesystem is cloned from @<snapshot>
	// The filesytems are those specified in SourceConfig, if any, otherwise this config.
	DeviceOrigin Pattern `yaml:"device-origin,omitempty"`

	// Experimental: source-config specifies a config file that is used to determine:
	//   - The LXD project, container, and snapshot to clone when launching the instance.
	//   - The source filesystems used for cloning filesystems or copying device directories.
	// The name of the instance used for the source filesystems
	// is the base name of the filename, without the extension.
	// Various parts of these items can be overriden by other source properties above
	SourceConfig HostPath `yaml:"source-config,omitempty"`
}

// OS specifies the container OS
type OS struct {
	// Name if the name of the container image, without the version number.
	// All included configuration files should have the same OS Name.
	// Supported OS names are "alpine", "debian", "ubuntu".
	// Support for an OS is the ability to determine the LXD image, install packages, create users, set passwords
	Name string `yaml:"name"`

	// Image is the image name (with optional remote), used when launching a container.
	// If Image is missing, an image name is constructed from Version.
	Image Pattern `yaml:"image"`
}

// Filesystem is a ZFS filesystem or a plain directory that is created when an instance is created
// The disk devices of an instance are created as subdirectories of a Filesystem
type Filesystem struct {
	// Pattern is a pattern that is used to produce the directory or zfs filesystem
	// If the pattern begins with '/', it is a directory
	// If it does not begin with '/', it is a zfs filesystem name
	Pattern Pattern `yaml:"pattern"`
	// Zfsproperties is a list of properties that are set when a zfs filesystem is created or cloned
	Zfsproperties map[string]string `yaml:"zfsproperties,omitempty"`
	// Destroy allows lxdops destroy the filesystem when requested.
	Destroy bool `yaml:"destroy,omitempty"`
	// Transient filesystems are not backed-up, exported, or imported.
	Transient bool `yaml:"transient,omitempty"`
}

// A Device is an LXD disk device that is attached to the instance profile, which in turn is attached to a container
type Device struct {
	// path is the device "path" in the LXD disk device
	Path string `yaml:"path"`

	// filesystem is the Filesystem Id that this device belongs to
	// If it is empty, dir should be an absolute path on the host
	Filesystem string `yaml:"filesystem"`

	// dir is the subdirectory of the Device, relative to its Filesystem
	// If empty, it default to the device Name
	// If dir == ".", the device source is the same as the Filesystem directory
	// Rarely used:
	// dir goes through pattern substitution, using parenthesized tokens, for example (instance)
	// dir may be absolute, but this is no longer necessary now that filesystems are specified, since one can define the "/" filesystem.
	Dir Pattern `yaml:"dir,omitempty"`

	// readonly - make the device readonly
	Readonly bool `yaml:"readonly,omitempty"`
}
