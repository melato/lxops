# Examples
These examples demonstrate the capabilities of lxops configuration files.
The commands to run them are described in the [Tutorial](md/tutorial.md)

## install packages
{{- $path := printf "%s/templates/ssh.yaml" .tutorial}}
This *lxops* configuration file can be used to install packages to a base alpine image,
which can then be published to an image containing the selected packages:

{{template "config.tpl" (Config.Args (printf "%s/templates/ssh.yaml" .tutorial) .).WithHeading "##"}}

{{- $build := printf "%s/templates/build.sh" .tutorial}}

## create a container with attached disk device
This *lxops* configuration file creates a container
with an external disk device, mounted at /home

It also creates a user in the container.

If we rebuild the container by deleting it and launching it again, the user will be created during each launch,
but the files in the user's home directory will be preserved across the rebuild.

{{- $path := (printf "%s/containers/ssh.yaml" .tutorial)}}
{{template "config.tpl" (Config.Args (printf "%s/containers/ssh.yaml" .tutorial) .).WithHeading "##"}}
