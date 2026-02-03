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

See [Filesystems]({{.lxops_doc_url}}/Filesystem.md)

