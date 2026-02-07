# Tutorial
## setup

First, create an lxops properties file:
```
mkdir -p ~/.config/lxops
touch ~/.config/lxops/properties.yaml
```

Initially we will need an *fshost* property with the name of a zfs filesystem to use
for creating instance-specific disk devices.

You can edit ~/.config/lxops/properties.yaml by hand:
```
fshost: tank/demo/host
```

or use:
```
lxops property set fshost tank/demo/host
```

Replace *tank/demo/host* with a valid empty zfs filesystem.

## create an alpine ssh image
We will start with a base alpine image and install packages to it,
using this lxops file:


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






To create the image, run the commands:
```
cd tutorial/templates

# launch a container called "ssh-template" and apply the cloud-config files
lxops launch -name ssh-template ssh.yaml

# if you script this, wait a few seconds to give some time
# to the container installation scripts to complete
sleep 5

# create a snapshot from this container
incus stop ssh-template
incus snapshot create ssh-template copy
incus publish ssh-template/copy --alias alpine-ssh

# list the new image
incus image list alpine-ssh

# delete the container.  We no longer need it.
lxops delete -name ssh-template ssh.yaml

```

We used a temporary container *ssh-template* and created image *alpine-ssh*.

## create ssh containers
We will create two ssh containers from this image.

We could create containers simply by using "incus launch":
```
incus launch alpine-ssh ssh
```

But we want to do more:
- create a user
- use a container-specific disk device for the /home directory,
so the user's home directory is stored independently from the guest OS.

```
cd tutorial/containers

## lxops extract copies files from the image to a
## "alpine-ssh" template container
## It copies the files that are needed to create
## instance-specific non-root disk devices
lxops extract -name alpine-ssh ssh.yaml

lxops launch -name ssh1 ssh.yaml
lxops launch -name ssh2 ssh.yaml

```

If in the future, we rebuild the ssh image and want to rebuild the containers that use it:
```
#!/bin/sh

incus stop ssh1
lxops delete -name ssh1 ssh.yaml
lxops launch -name ssh1 ssh.yaml

```

Rebuilding does not preserve container ip addresses.

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
an lxops propertiy, and it is replaced by the corresponding property value.

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





