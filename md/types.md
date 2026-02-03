# lxops configuration types


- [Config](Config.md)
- [Pattern](Pattern.md)
- [HostPath](HostPath.md)
- [Filesystem](Filesystem.md)
- [Device](Device.md)

# Conditional configuration
*include* and *cloud-config-files* paths in lxops files go through variable substitution
and are used only if all referenced variables exist.

The characters "|;,:" are reserved.  They are not allowed in these path names.
