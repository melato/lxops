lxops is a Go program that launches and configures LXD or Incus containers.

# Main Features

## launching containers
- lxops launches containers using lxops configuration files

Each configuration file specifies, among other things:
- An image to launch the container from.  Alternately, the container can be cloned from a snapshot of another container.
- A list of profiles to apply to the container
- A list of standard [cloud-init](https://cloudinit.readthedocs.io) files that configure the container.
The configuration is done via the instance server API.  It does not use cloud-init packages in the container.
- A list of ZFS disk devices to create and add to the container.
By using lxops, you can create, snapshot, delete, or restore these devices with a single command per container.

A configuration file can include other configuration files.

## rebuilding containers
"lxops rebuild" reconstructs a container by attaching a new image to a container's attached devices.

- stops the container
- deletes the container
- launches and configures the container again
- It reuses any attached devices that it manages for the containers

Rebuilding tries to preserve the container's IP address and profiles.  This can be disabled by command-line flags.

The goal is to be able to replace the root device of a container with an updated one,
without disrupting the applications in the container, except for a reboot.
To do this correctly requires knowledge of the files that the applications of interest use,
so that any changes to these files are placed in external filesystems
or they are reconfigured when the container is launched or rebuilt.

I typically attach disk devices to:
- /log
- /tmp
- /home
- /etc/opt
- /var/opt
- /opt
- /usr/local/bin



# Examples
For detailed examples, see the separate [lxops.demo](https://github.com/lxops.demo) repository.

# Compile
This repository does not produce an executable.
It depends on an external implementation of an Instance API.
There are two other repositories that implement this API and produce executables:
- lxops_lxd produces lxops-lxd which works with LXD
- lxops_incus produces lxops-incus which works with Incus

Once you build one of these two executables, rename it to "lxops".


To build lxops_lxd:
	cd lxops_lxd/main
	git log -1 --format=%cd > version
	# or: date > version
	go build lxops-lxd.go
	ln -s lxops-lxd lxops

# Use Case
You want to have a set of containers that have the same guest OS, packages, users, etc.
in order to run websites, for example a container that has a web server, PHP, and several libraries.
You want to have one container per website, but avoid installing and upgrading packages for each container.

You can create a template container with the packages that you want, snapshot it, and then clone it for each website.
You can do this, using lxc copy:
```
lxc copy <template-container>/<snapshot> <working-container>
```
Assuming that you use ZFS or another copy-on-write filesystem, the root filesystem of each working container is a clone of the template container root filesystem, so it uses very little disk space.  Creating the working containers is a relatively fast operation, that does not download and install packages.

But what do you do to upgrade the containers?  If you upgrade each one separately, the working container root filesystems start to diverge from their template.  In addition, you will be downloading and installing the same upgrades multiple times (possibly tens or hundreds of times per LXD host).

lxops facilitates the following strategy:
- You structure each working container so that application data is not on the root filesystem, but on external disk devices.  Therefore you can replace the root filesystem with a new one, without losing the application data.
- To upgrade, you first create a new template container with the upgrade you want.  Alternately, you can upgrade the existing template container and create a new template snapshot.

Then, for each working container:
- Delete the container
- Clone the container from the template
- re-attach the existing external disk devices
- start the container.  The container should now run with its new OS and the old application data.
This process takes a few seconds per container, during which the container will be offline.

This is what lxops is for.

# lxops config file
An lxops config file is a yaml file that provides instructions about how to build
a template container or a working container.  It can include other config files.
Detailed documentation is in the Go docs:
	cd lxops/cfg
	go doc Config

It provides:
- The image or container/snapshot to create the container from.
- Filesystems and external disk devices to create or use for the container.
- LXD profiles to attach to the container
- A list of cloud-config files to use for configuring the container.

# Filesystems/Disk Devices
lxops can create 0 or more zfs filesystems for each container.
Each filesystem is parameterized by container name, so each container is automatically assigned its own private external filesystems.
If they do not exist, they are created.  If they exist, they are used as is.

The external disk devices that lxops manages are subdirectories of these filesystems. 
I typically use one filesystem for /var/log, one filesystem for /tmp,
and one filesystem for /var/opt, /etc/opt, /opt, /home, /usr/local/bin.
If they do not already exist, they can be copied from the corresponding devices of a template container.

# cloud-config files
lxops uses a subset of the cloud-config file format to configure containers internally.
The cloud-config files are applied directly using the LXD API,
without requiring that the container supports cloud-init.

The cloud-config modules (sections) that are supported are:
- packages
- write_files
- users
- runcmd

For more details, see:
- https://cloudinit.readthedocs.io/en/latest/reference/examples.html
- https://github.com/melato/cloudconfig
- https://github.com/melato/cloudconfiglxd

# Goals
- Separate container OS from application data, so that application data
is not on the container root filesystem.
- Create containers by cloning the root filesystem of a template container
(using lxc copy of a snapshot), and creating/copying, or cloninig additional filesystems.

It reads container and filesystem configuration from YAML files.
Configuration inside the container is done via yaml files in the cloud-init format (#cloud-config).

# Description

lxops launches, copies, and deletes *lxops instances*.

An lxops **instance** is:
- An LXD container
- A set of ZFS filesystems specific to this container
- A set of disk devices that are in these filesystems and are attached to the container (via a profile)
- An LXD profile that specifies the attached disk devices

A Yaml instance configuration file specifies how to launch and configure an instance.  It specifies:
- How to launch the container, from an image, or by copying a snapshot of another container
- ZFS Filesystems to create and zfs properties for these filesystems
- How to create these filesystems:  create them from scratch, or copy them (rsync) or clone them from another instance.
- Disk Devices in these filesystems that are attached to the container
- An LXD profile to create with the instance devices
- Additional LXD profiles to attach to the container
- cloud-config files to apply to the container

Several configuration elements can be parameterized with properties such as the instance name, project, and user-defined properties.

More detailed documentation of configuration elements is in the file Config.go

# Project Support

lxops has some limited support for LXD/Incus projects and can clone instances across projects.
I have no good use case for it and I don't use it, so should be considered unsupported feature.
I find it simpler to keep all my instances in a single project.

By default, lxops will use the current project, as detected by looking at the lxc/incus user config files.

# External Programs

lxops calls these external programs, on the host, with *sudo* when necessary:
- lxc or incus (It mostly uses the InstanceServer API, but uses the "lxc" or "incus" commands for certain operations)
- zfs
- rsync
- chown
- mkdir
- mv

lxops calls these external programs, in the container, via cloud-config files:
- sh
- chpasswd
- chown
- OS-specific commands for adding packages and creating users

