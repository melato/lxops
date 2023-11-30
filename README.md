lxops launches LXD or Incus containers from config (lxops) files that 
specify how the instance should be launched and configured.

# Goals
- launch instances from a config file
- apply cloud-config files to an instance.
- create instance-specific filesystems when creating an instance,
and attach disk devices from these filesystems.
- rebuild an instance from a new image, while preserving its filesystems and configuration

Configuration inside a container is specified using standard [cloud-config](https://cloud-init.io/) files,
which are applied using the InstanceServer API, without using any cloud-init packages.
See [github.com/melato/cloudconfig](https://github.com/melato/cloudconfig) for what is supported.

# Examples
You can find examples in the separate [lxops.examples](https://github.com/melato/lxops.examples) repository.

Here is a simple configuration file (example.yaml):

```
#lxops-v1
ostype: alpine
image: images:alpine/3.18
profiles:
- default
cloud-config-files:  
- ../packages/base.cfg
- ../packages/bash.cfg
- ../cfg/doas.cfg
- ../cfg/user.cfg
```

You can create containers a1, a2, using these commands:
```
	lxops launch -name a1 example.yaml
	lxops launch -name a2 example.yaml
```
It's even better if you create an image from this configuration, and create the containers from your image.
The examples repository demonstrates that.

# Compile
This project is provided as Go source code only at this point.
It provides code for two separate executables, one for LXD and one for Incus.
Compile either one and create an "lxops" link to the one that you want to use.

To compile lxops-lxd:
```
cd ./impl/lxd/main
go install lxops-lxd.go
```

To compile lxops-incus:
```
cd ./impl/incus/main
go install lxops-incus.go
```

The code is split into three different Go modules, so each module has its own dependencies.

# Disk Devices
A central feature of lxops is the ability to create and attach external disk devices to a container it launches.

The intent is that the combination of external devices and configuration
makes it possible to rebuild a container with a new image without losing data.

ZFS and plain directory devices are supported.  ZFS is the implementation tested most.

For ZFS devices, lxops will:
- create ZFS filesystems with names that include the container name.
- create directories in these filesystems, one for each device.
- change the permission of these directories so that they appear as owned by root:root in the container.
- copy files in these directories from template directories, if a device template is specified.
- add these directories as disk devices to the lxops profile associated with the container.

I typically attach disk devices to all these directories:
- /home
- /etc/opt
- /var/opt
- /opt
- /usr/local/bin
- /var/log
- /tmp

And make sure I put all application data in one of these directories (except /tmp, of course).

When I replace the root filesystem with a new image, my data persists.

An image may already have /log populated with files and directories, without which the image might not function properly.
For this reason, the configuration file has a *device-template* field that specifies another container whose devices will
be copied to the current container during lxops launch, using rsync.

## lxops files
An lxops file is a yaml configuration file that specifies how to launch an instance.
Documentation for lxops files is provided by the "lxops help config".

# commands
Here are the basic commands for managing a container with lxops. 

-name <container-name> can be omitted if it is the same as the config-file base name,
so these commands are equivalent:
	lxops launch myconfig.yaml
	lxops launch -name myconfig myconfig.yaml

## launch
	lxops launch -name <container> <config-file.yaml>
	
- creates any zfs filesystems that are specified in the config file.
  The filesystem names typically include the container name, so they are unique for each container.
- creates subdirectories in these filesystems
- copies files from template filesystems, if specified
- adds the subdirectories as disk devices to a new profile named <container>.lxops
- creates a container named with the profiles specified in the config file and the <container>.lxops profile
- configures the container by running the specified cloud-config files

If the device directories already exist, they are used without overwriting them.

## delete
	lxops delete -name <container> <config-file.yaml>

- deletes the container (it must already be stopped)
- deletes the container profile (<container>.lxops)

## rebuild
	lxops rebuild -name <container> <config-file.yaml>
	
- stops the container
- runs "lxc rebuild" or "incus rebuild" with the image specified in the config file
- starts the container
- applies the cloud-config files

Since the device directories exist, they are are preserved across the rebuild.
The result is that the container has a new guest OS, but runs with the old data.
For this to work, the container must be configured properly, via the cloud-config files.
	
## destroy
	lxops destroy -name <container> <config-file.yaml>

- deletes the container
- deletes the container profile (<container>.lxops)
- destroys the container filesystems

# External Programs

lxops calls these external programs, on the host, with *sudo* when necessary:
- lxc or incus (It mostly the InstanceServer API, but uses the "lxc" or "incus" commands for certain complex operations, like launch or rebuild).
- zfs
- rsync
- chown
- mkdir
- mv

lxops calls these container executables, via cloud-config files:
- /bin/sh
- chpasswd
- chown
- OS-specific commands for adding packages and creating users

# Stability/History
lxops is a continuation of lxdops,
which I have been using for years to manage containers with numerous configuration files.

Therefore, I have a personal interest to keep it working and maintain backward compatibility with the lxops configuration file format.  Nevertheless, I want to simplify it, so changes may happen.  Some half-baked features may be removed.

## configuration file format
lxops supports multiple configuration file formats.  There are currently two supported formats:
  - lxops-v1 - The latest format.
  - lxdops - An older format that lxdops used

backward compatibility is maintained by using migrators that convert a format to a newer format.
- "lxdops" files are converted to "lxops-v1" files.
- If and when there is a lxops-v2 format, there could be an lxops-v1 migrator that converts lxops-v1 files to lxops-v2 files.

Format migrators are chained, so lxdops files will be converted to lxops-v1 files and then to lxops-v2 files.
Therefore, all previous formats should be supported.

Format migrators convert bytes to bytes, so they do not need to depend on lxops data types.

You can use your own format, if you want.  You just have to write a format migrator and install it in your own main().

## command-line interface
lxops CLI changed somewhat from lxdops.
The main container management operations (launch, delete, destroy, rebuild) remain the same.

Recent changes:
- The former "instance" subcommands have been renamed to "i", to avoid confusion with LXD/Incus instances.
  These are informational/debugging commands for lxops configuration files ("instance" files) and should not be scripted.
- The "container" subcommands are renamed to "instance", to reflect LXD/Incus terminology.
  These are mostly informational commands for LXD/Incus instances that are not central to lxops and were not meant to be scripted.
