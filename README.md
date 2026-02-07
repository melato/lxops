lxops is a tool that assists launching Incus containers.

# Goals
It provides the following facilities on top of incus:

## scriptable launch and configure
Launch and configure instances from *lxops* yaml configuration files
and [cloud-config](https://cloud-init.io/) files.

## instance-specific disk devices
Create instance-specific filesystems when launching an instance,
and attach disk devices from these filesystems to the instance.
This has been designed and tested for zfs filesystems.

This allows separating the OS from your data and applications,
and upgrading a container by replacing the OS with a new image.

To do this, you must store application data to non-root disk devices.

You can store configuration files in non-root disk devices,
or re-generate them when re-launching the container, or a combination of the two.

The [lxops.examples](https://github.com/melato/lxops.examples) repository
shows how to do this for a few applications.

# Examples
These examples demonstrate the capabilities of lxops configuration files.
The commands to run them are described in the [Tutorial](md/tutorial.md)

## install packages
This *lxops* configuration file can be used to install packages to a base alpine image,
which can then be published to an image containing the selected packages:


**./tutorial/templates/ssh.yaml**
```
#lxops-v1
ostype: alpine
image: images:alpine/3.23
cloud-config-files:
- ../packages/ssh.cfg

```

It launches a container with the image *images:alpine/3.23*.


### cloud-config-files
After the container is launched and started,
the *cloud-config-files* are applied to it:

**tutorial/packages/ssh.cfg**
```
#cloud-config
packages:
# dhcpcd is needed for getting an ipv6 address via DHCP
- dhcpcd
- openssh

```






## create a container with attached disk device
This *lxops* configuration file creates a container
with an external disk device, mounted at /home

It also creates a user in the container.

If we rebuild the container by deleting it and launching it again, the user will be created during each launch,
but the files in the user's home directory will be preserved across the rebuild.

**./tutorial/containers/ssh.yaml**
```
#lxops-v1
ostype: alpine
image: alpine-ssh
device-template: alpine-ssh
include:
- include/device.yaml
cloud-config-files:
- include/user.cfg

```

The lxops file includes other lxops files listed in **include**.
After includes are merged with the current file, the resulting configuration is:
```
ostype: alpine
image: alpine-ssh
device-template: alpine-ssh
profile-pattern: (instance).lxops
device-owner: 1000000:1000000
filesystems:
  host:
    pattern: (fshost)/(instance)
    destroy: true
devices:
  home:
    path: /home
    filesystem: host
cloud-config-files:
- tutorial/containers/include/user.cfg

```

It launches a container with the image *alpine-ssh*.


### filesystems

The **pattern** field of each filesystem, specifies the host location of the filesystem.

It first goes through variable substitution:
- (instance) is replaced by the instance name.
- Any other parenthesized expression is treated as the name of
an lxops property, and it is replaced by the corresponding property value.

If the resulting $pattern does not begin with '/', it is treated as a zfs filesystem,
and it is created using
```
   sudo zfs create $pattern
```

If $pattern begins with '/', it is treated as a plain directory,
and it is created using
```
   sudo mkdir -p $pattern
```

### disk devices

Each filesystem can have multiple disk devices.

Each device is a subdirectory of its corresponding filesystem,
except for a device with *dir* set to '.', which is mapped to the filesystem directory itself,
without using a subdirectory.

If the device directory of a disk device already exists, it is left alone.


### device-template

If device-template is specified:
```
device-template: *alpine-ssh*
```

The device is created with
```
    sudo mkdir -p $device_dir
```

Its content is initialized using:
```
   rsync -av $template_device_dir/ $device_dir/
```

- $device_dir is the full path of the device,
composed from the filesystem directory and the device subdirectory

- $template_device_dir is composed the same, way, except that the **(instance)**
variable in the filesystem pattern is replaced by the value of the *device-template*
field (*alpine-ssh*) in the lxops config file.

### device-template initialization

The *device-template* disk device directories can be initialized from the image,
using:
```
  lxops extract -name alpine-ssh ./tutorial/containers/ssh.yaml
```



The **device-owner** field indicates the uid and gid of the root user in the container,
as seen by the host.

*lxops extract* copies the device files from the image, using *rsync*, and then
modifies the ownership of each file and directory by adding the *device-owner* values.

After the filesystems are created and the disk device directories populated,
a profile is created for the instance, using the profile name specified in the *profile-pattern* field.
This profile specifies the disk devices and is attached to the new instance.



### cloud-config-files
After the container is launched and started,
the *cloud-config-files* are applied to it:

**tutorial/containers/include/user.cfg**
```
#cloud-config
users:
- name: demo
  uid: 1000
  sudo: true
  groups: wheel,adm

```







# [Further Documentation](md/index.md)
