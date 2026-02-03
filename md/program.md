# External Programs

lxops calls these external programs, on the host, with *sudo* when necessary:
- incus or lxc (It mostly uses the InstanceServer API, but uses the "incus" or "lxc" commands for certain complex operations, like launch).
- zfs
- rsync
- chown
- mkdir
- mv

lxops calls these container executables, via cloud-config files:
- /bin/sh
- chpasswd
- chown
- OS-specific commands for adding packages and creating users

# Stability/History
lxops is a continuation of lxdops.
I have been using it for years, since 2020 to manage containers, using dozens of configuration files.

Therefore, I have a personal interest to keep it working and maintain backward compatibility with the lxops configuration file format.
Nevertheless, changes may happen.  Some half-baked features may be removed.

