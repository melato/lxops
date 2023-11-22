#### data types
Each field below has its Go data type in parenthesis.
Custom data types are:
- Pattern: a string that is subject to property substitution.
The pattern (x) is replaced with the value of the lxops property "x".
The built-in property "instance" is the name of the instance.

- HostPath: a string that denotes a file path, either absolute,
or relative to the config file directory.

## ostype (string)
OS type for cloudconfig.  Two ostypes are provided:
- alpine, for alpine distributions.  Uses apk commands to add packages.

- debian, for debian-based distributions.  Uses apt-get commands to add packages.

## image (Pattern)
The image name (with optional remote), used for launching a container.

## profiles ([]string)
The profiles that an instance will be launched with,
excluding the lxops profile which is always added.

When rebuilding an instance, the instance is relaunched either with these profiles, or with the profiles that the instance previously had (which may have changed since the previous launch).

While the instance is being created and configured it has these profiles
minus ProfilesRun plus ProfilesConfig

## cloud-config-files	([]string)
A list of cloud-config files to apply to an instance when launching or rebuilding.
The files can be specified using absolute paths or paths relative to the enclosing configuration file.

## include ([]string)
Include is a list of other configs that are to be merged with this config.
Include paths are either absolute or relative to the path of the including config.

TODO:  Clarify merge rules.

## filesystems (map[string]*Filesystem)
Filesystems are zfs filesystems or plain directories that are created
when an instance is created.  Devices are created inside filesystems.

The keys in the filesystems map are used ion the devices section to reference each filesystem.

Each Filesystem has:
### pattern (Pattern)
A pattern that is used to derive the filesystem path on the host.

Example: (fsroot)/log/(instance)

Filesystems should typically include the "instance" property so that they are unique for each instance.

If the resulting path begins with "/", it is a plain directory.

If it does not begin with "/", it is a zfs filesystem.

### destroy (bool)
Specifies that the filesystem should be destroyed by "lxops destroy". 

This should be set to true in most cases.
It is a safety feature to remind you that this filesystem may be destroyed by lxops.

### transient (bool)
Marks a filesystem as transient.

Transient filesystems are not exported or imported and do not need to be backed up.

A device where a transient filesystem is appropriate is /tmp.

### zfsproperties (map[string]string)
A list of zfs properties to set when creating or cloning a zfs filesystem.
They are passed to zfs using the -o flag.

## devices (map[string]*Device)
Devices are disk devices that are directories within the instance filesystems

They are created and attached to the container via the instance profile.

Each Device has:
### path (string)
The path of the device in the instance.

### filesystem (string)
The key of the filesystem that the device is in.

### dir (string)
A subdirectory for the device in its filesystem.

If dir is missing, it defaults to the key of the device in the devices map.

If dir is ".", the device host directory is the directory of the filesystem, so there is no sub-directory created.
	

## device-owner (Pattern)
The owner (uid:gid) to set when creating devices.
It is typically 1000000:1000000,
because LXD and Incus usually map host uid/gid 1000000 to instance uid/gid 0.

## profile-pattern (Pattern)
Specifies how the lxops instance profile should be named.
Defaults to "(instance).lxops"

## properties (map[string]string)
Properties provide key-value pairs used for pattern substitution.

They override built-in or global properties

Properties from all included files are merged before they are applied.

Properties cannot override non-empty properties,
in order to avoid unexpected behavior that depends on the order of included files.

## origin (Pattern)
origin is the name of a container and a snapshot to clone from.
It has the form [<project>_]<container>[/<snapshot>]
It overrides SourceConfig

## device-template (Pattern)
device-template is the name of an instance, whose devices are copied (using rsync)
to a new instance when launching.  The instance does not need to exist.

The device host directories are calculated using the current config files, using the substituted device-template as an instance name.

## device-origin (Pattern)
device-origin is used to specify that the instance filesystems should be cloned (using zfs clone) from a snapshot of the filesystems of another instance.

It should have the form <instance>@<snapshot>, where @<snapshot> is the short zfs snapshotname of the instance filesystems.

## stop (bool)
Stop specifies that the container should be stopped at the end of the configuration.
Use when creating images or container snapshots that will be used to clone other containers.

## snapshot (string)
Snapshot specifies that that the container should be snapshoted with this name at the end of the configuration process.

# Experimental Properties
Any config fields not listed here are experimental and may be removed.
