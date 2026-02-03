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

To do this, you must have arranged for all application configuration and data to
either be stored in non-root disk devices, or re-generated when re-launching the container.

The [lxops.examples](https://github.com/melato/lxops.examples) repository
shows how to do this for a few applications.

{{template "tutorial.tpl" .}}

{{template "compile.tpl" .}}

# [Further Documentation](md/index.md)
