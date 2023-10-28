lxops launches and configures LXD or Incus containers,
using configuration files that have instructions about how to build and configure a container.
It can also create and attach disk devices to these containers.

# Example
You can find examples in the separate [lxops.examples](https://github.com/melato/lxops.examples) repository.

A simple example of an lxops configuration file is:

[(alpine/containers/example.yaml)](https://github.com/melato/lxops.examples/blob/main/alpine/containers/example.yaml)
```
#lxops-v1
ostype: alpine
image: images:alpine/3.18
description: launch container from images:
profiles:
- default
cloud-config-files:  
- ../cfg/doas.cfg
- ../cfg/dhcpcd.cfg
- ../cfg/interfaces.cfg
```

For getting started, you can delete the cloud-config-files section,
so the configuration file can be used without any other files.
Or run them in the lxops.examples project, so you have all the dependencies.

lxops can use this file above as follows:

## launch
	lxops launch -name example1 example.yaml

This launches a container called "example1".
It applies the specified profiles (in this case just "default"),
plus an additional lxops profile that it creates just for this container (e.g. example1.lxops).
It also runs the specified #cloud-config files inside the container.

The cloud-config files are applied using the InstanceServer API.  It does not use any cloud-init packages.
See [github.com/melato/cloudconfig](https://github.com/melato/cloudconfig) for what is supported.


## delete
	lxops delete -name example1 example.yaml

This deletes the container and its lxops profile.
The container must already be stopped.

## rebuild
	lxops rebuild -name example1 example.yaml
	
	
Rebuilding stops the container, deletes it, and launches it again.

This is useful if the container is configured to keep persistent configuration and data in
attached disk devices that are preserved during the rebuilding.

lxops provides the ability to manage attached disk devices.  This is not shown in the basic example above.

This way, you can replace the container guest OS, without losing your data.

Rebuilding preserves the old container's profiles and IP addresses.  This behavior can be disabled by flags.

## destroy
	lxops destroy -name example1 example.yaml

In this example, destroy is the same as delete.

More advanced configuration files can specify disk devices that lxops create during launch and destroys with the "destroy" command.

# Disk Devices
A central feature of lxops is the ability to create and attach external disk devices to a container it launches.

The intent is that the combination of external devices and configuration
makes it possible to rebuild a container with a new image without losing data.

ZFS and plain directory devices are supported.  ZFS is the implementation tested most.

For ZFS devices, lxops will:
- create 0 or more ZFS filesystems with names that include the container name.
- create directories in these filesystems, one for each device.
- change the permission of these directories so that they appear as owned by root:root in the container.
- optionally copy files in these directories from template directories.
- add these directories as disk devices to the lxops profile associated with the container.

I typically attach disk devices to all these directories:
- /var/log
- /tmp
- /home
- /etc/opt
- /var/opt
- /opt
- /usr/local/bin

An image may already have /log populated with files and directories, without which the image might not function properly.
For this reason, the configuration file has a *device-template* field that specifies another container whose devices will
be copied to the current container during lxops launch, using rsync.

# Compile
This project is a Go library that communicates with an instance server via a backend interface.
It does not produce an executable.

It is used by two other projects that implement this interface for LXD or for Incus:
- [lxops_lxd](https://github.com/melato/lxops_lxd)
- [lxops_incus](https://github.com/melato/lxops_incus)

Once you build one of these two executables, rename it to "lxops" (or link "lxops" to it).

To build lxops_lxd:
```
    git clone https://github.com/melato/lxops_lxd.git
    cd lxops_lxd/main
    git log -1 --format=%cd > version
    # or: date > version
    go build lxops-lxd.go
    ln -s lxops-lxd lxops
```

# External Programs

lxops calls these external programs, on the host, with *sudo* when necessary:
- lxc or incus (It mostly uses the InstanceServer API, but uses the "lxc" or "incus" commands for certain complex operations)
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
lxops supports multiple configuration file formats.  There are currently two formats:
  - lxdops - the same format that lxdops used
  - lxops-v1 - a simpler format with os replaced with ostype, image.

backward compatibility is maintained by using migrators that convert a format to a newer format.
- "lxdops" files are converted to "lxops-v1" files.
- when there is lxops-v2 format, there will be an lxops-v1 migrator that converts lxops-v1 files to lxops-v2 files.

Format migrators are chained, so lxdops files will be converted to lxops-v1 files and then to lxops-v2 files.
Therefore, all previous formats should be supported.

Format migrators convert bytes to bytes, so they do not need to depend on lxops data types.

Some configuration fields may be removed, because they are rarely used or not fully implemented.
If this happens, a new format migrator may throw an error if it encounteres an obsolete field, so that you know
that you need to fix the failed file.

You can use your own format, if you want.  You just have to write a format migrator and install it in your own main().

## command-line interface
lxops CLI changed somewhat from lxdops.
The main container management operations (launch, delete, destroy, rebuild) remain the same.

Recent changes:
- The former "instance" subcommands have been renamed to "i", to avoid confusion with LXD/Incus instances.
  These are informational/debugging commands for lxops configuration files ("instance" files) and should not be scripted.
- The "container" subcommands are renamed to "instance", to reflect LXD/Incus terminology.
  These are mostly informational commands for LXD/Incus instances that are not central to lxops and were not meant to be scripted.

## other differences from lxdops:
- lxops uses no LXD code. All such dependencies are in a separate repository (lxops_lxd).
- lxops uses the configuration directory ~/.config/lxops/ instead of ~/.config/lxdops/
- The default profile suffix is ".lxops" instead of ".lxdops".  This can be changed in configuration files.


