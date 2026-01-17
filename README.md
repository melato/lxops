lxops launches Incus containers from config (lxops) files that 
specify how the container should be launched and configured.

# Goals
## scriptable launch and configure
Launch and configure instances from yaml configuration files and cloud-config files.

Configuration inside a container is specified using
standard [cloud-config](https://cloud-init.io/) files,
which are applied using the Incus API, without using any cloud-init packages.
See [github.com/melato/cloudconfig](https://github.com/melato/cloudconfig) for what is supported.

## instance-specific disk devices
create instance-specific filesystems when creating an instance,
and attach disk devices from these filesystems to the instance.

This allows upgrading a container by rebuilding it with a new image,
if you have arranged for all application configuration and data to
be in non-root disk devices.
The [lxops.examples](https://github.com/melato/lxops.examples) repository
shows how to do this for a few applications.

# Examples
You can find working examples in the separate [lxops.examples](https://github.com/melato/lxops.examples) repository.

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

I put all application data and configuration in one of these directories (except /tmp).

When I replace the root filesystem with a new image, my data persists.

## lxops files
An lxops file is a yaml configuration file that specifies how to launch an instance.
Documentation for lxops files is provided by the "lxops help config".

# commands
Here are the basic commands for managing a container with lxops. 

-name <container-name> can be omitted if it is the same as the config-file base name,
so these commands are equivalent:
```
lxops launch myconfig.yaml
lxops launch -name myconfig myconfig.yaml
```

## extract
See the examples repository for more detailed examples of using *extract* and *launch*.

```
lxops extract -name <name> <config-file.yaml>
```

*extract* creates filesystems and disk device directories for a container.
It copies files from an image to these device directories.


## launch
```
lxops launch -name <container> <config-file.yaml>
```
	
- creates any zfs filesystems that are specified in the config file.
  The filesystem names typically include the container name, so they are unique for each container.
- creates subdirectories in these filesystems
- copies files from template filesystems, if specified.
  The template filesystems can be copied from the container's image by using the *extract* command.
- adds the subdirectories as disk devices to a new profile named <container>.lxops
- creates a container with the profiles specified in the config file and the <container>.lxops profile
- configures the container by running the specified cloud-config files

If the filesystems or device directories already exist, they are used without overwriting them.

## delete
```
incus stop <container>
lxops delete -name <container> <config-file.yaml>
```

- deletes the container (it must already be stopped)
- deletes the container profile (<container>.lxops)
- does not touch the container filesystems
- *delete* will do nothing if the container is running

## rebuild
```
incus stop <container>
lxops delete -name <container> <config-file.yaml>
lxops launch -name <container> <config-file.yaml>
```

There is an *lxops rebuild* command that is supposed to do the equivalent,
but it is not adequately tested.  I recommend running the three commands above.

Since the device directories exist, they are are preserved across the rebuild.
The result is that the container has a new guest OS, but runs with the old data.
For this to work, the container must be configured properly, via the cloud-config files.
	
## destroy
```
lxops destroy -name <container> <config-file.yaml>
```

Does everything that "delete" does and also destroys the container filesystems.
- *destroy* will do nothing if the container is running

run *lxops help filesystem* for more information about filesystems.

# External Programs

lxops calls these external programs, on the host, with *sudo* when necessary:
- incus or lxc (It mostly uses the InstanceServer API, but uses the "incus" or "lxc" commands for certain complex operations, like launch).
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
lxops is a continuation of lxdops.
I have been using it for years, since 2020 to manage containers, using dozens of configuration files.

Therefore, I have a personal interest to keep it working and maintain backward compatibility with the lxops configuration file format.
Nevertheless, changes may happen.  Some half-baked features may be removed.

## configuration file format
lxops supports multiple configuration file formats, via configuration file migrators.
This mechanism was used to maintain backward-compatibility when the format changed.

see *lxops help topics config-format*

# Compile
This project is provided only as Go source code at this point.

After compiling, link the resulting executable (lxops-incus or lxops-lxd) to "lxops" and put "lxops" in your path.
Various examples and scripts execute "lxops".

## get the code
```
git clone github.com/melato/lxops
cd lxops
```

## Compile for Incus

```
cd ./impl/incus/main
go install lxops-incus.go
```

## Compile for LXD

```
cd ./impl/lxd/main
go install lxops-lxd.go
```

I have stopped using lxops-lxd, so I don't know if it still works.
