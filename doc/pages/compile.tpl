# Compile
This project is provided only as Go source code at this point.

After compiling, link the resulting executable (lxops-incus or lxops-lxd) to "lxops" and put "lxops" in your path.
Various examples and scripts execute "lxops".

## get the code
```
git clone github.com/melato/lxops
cd lxops
```

## Compile for Incus

```
cd ./impl/incus/main
go install lxops-incus.go
```

## Compile for LXD

```
cd ./impl/lxd/main
go install lxops-lxd.go
```

lxops can also compile with LXD instead of incus, but
I have stopped using and testing the LXD version.
Some recent features were added only to the Incus version.

When I migrated from LXD to incus, I did not actually migrate any containers.
I rebuilt my containers with Incus, as I was already doing with LXD.
I used the same lxops configuration files to build new Incus images and containers,
and attached the existing disk devices to the new containers.
The disk devices and their contents are not tied to either Incus or LXD.
They are simply zfs filesystems and directories.

