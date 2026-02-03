
# HostPath
```
HostPath is a string that represents a file path relative to a
parent path.
If it is absolute, it is used as is.
Otherwise, it is joined with the path of its parent.

The parent path depends on the context.

For config "include" or "cloud-config-files" files,
the parent path is the path of the enclosing config.

For a device dir, the parent path is the path of its filesystem.

Conditional paths:
"include" or "cloud-config-files" paths go through variable substitution.
If a path variable is not defined, the path is not used.
Variable substitution in paths uses global and command-line properties only.
It does not use properties defined in config files.
```
