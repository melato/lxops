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
{{- $path := printf "%s/templates/ssh.yaml" .tutorial}}
We will start with a base alpine image and install packages to it,
using this lxops file:

{{template "config.tpl" (Config.Args (printf "%s/templates/ssh.yaml" .tutorial) .).WithHeading "##"}}

{{- $build := printf "%s/templates/build.sh" .tutorial}}

To create the image, run the commands:
```
cd tutorial/templates
{{file (printf "%s/templates/build.sh" .tutorial)}}
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
{{file "tutorial/containers/launch.sh"}}
```

If in the future, we rebuild the ssh image and want to rebuild the containers that use it:
```
{{file "tutorial/containers/rebuild.sh"}}
```

Rebuilding does not preserve container ip addresses.

{{- $path := (printf "%s/containers/ssh.yaml" .tutorial)}}
{{template "config.tpl" (Config.Args (printf "%s/containers/ssh.yaml" .tutorial) .).WithHeading "##"}}
