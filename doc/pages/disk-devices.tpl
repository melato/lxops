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

