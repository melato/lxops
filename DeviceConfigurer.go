package lxops

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"melato.org/lxops/cfg"
	"melato.org/lxops/srv"
	"melato.org/lxops/util"
	"melato.org/script"
)

type DeviceConfigurer struct {
	Config *cfg.Config
	// NoRsync - do not rsync devices.  Use when importing.
	NoRsync bool
	Trace   bool
	DryRun  bool
}

func NewDeviceConfigurer(config *cfg.Config) *DeviceConfigurer {
	t := &DeviceConfigurer{Config: config}
	return t
}

func (t *DeviceConfigurer) NewScript() *script.Script {
	return &script.Script{Trace: t.Trace, DryRun: t.DryRun}
}

func (t *DeviceConfigurer) chownDir(scr *script.Script, dir string, owner string) {
	if owner != "" {
		scr.Run("sudo", "chown", owner, dir)
	}
}

func (t *DeviceConfigurer) CreateDir(dir string, owner string) error {
	if !util.DirExists(dir) {
		script := t.NewScript()
		script.Run("sudo", "mkdir", "-p", dir)
		//err = os.Mkdir(dir, 0755)
		t.chownDir(script, dir, owner)
		return script.Error()
	}
	return nil
}

func (t *DeviceConfigurer) CreateFilesystem(fs *InstanceFS, originDataset string, originfs *InstanceFS, owner string) error {
	if fs.IsDir() {
		fs.IsNew = true
		return t.CreateDir(fs.Dir(), "")
	}

	doClone := originDataset != "" && !originfs.Filesystem.Transient
	var args []string
	if doClone {
		args = []string{"zfs", "clone", "-p"}
	} else {
		args = []string{"zfs", "create", "-p"}

	}

	// add properties
	for key, value := range fs.Filesystem.Zfsproperties {
		args = append(args, "-o", key+"="+value)
	}

	if doClone {
		args = append(args, originDataset)
	}
	args = append(args, fs.Path)
	s := t.NewScript()
	s.Run("sudo", args...)
	if originDataset == "" {
		t.chownDir(s, fs.Dir(), owner)
		fs.IsNew = true
	}
	return s.Error()
}

func (t *DeviceConfigurer) CreateFilesystems(instance, origin *Instance, snapshot string, owner string) error {
	paths, err := instance.Filesystems()
	if err != nil {
		return err
	}
	var originPaths map[string]*InstanceFS
	if origin != nil {
		originPaths, err = origin.Filesystems()
		if err != nil {
			return err
		}
		for id, path := range paths {
			if !path.IsZfs() {
				return errors.New("cannot use origin with non-zfs filesystem: " + id)
			}
		}
	}
	var pathList []*InstanceFS
	for _, path := range paths {
		if origin != nil || !util.DirExists(path.Dir()) {
			pathList = append(pathList, path)
		}
	}
	InstanceFSList(pathList).Sort()

	for _, path := range pathList {
		var originDataset string
		var originPath *InstanceFS
		if path.IsZfs() {
			var exists bool
			originPath, exists = originPaths[path.Id]
			if exists {
				originDataset = originPath.Path + "@" + snapshot
			}
		}
		err := t.CreateFilesystem(path, originDataset, originPath, owner)
		if err != nil {
			return err
		}
	}
	return nil
}

func parseOwner(owner string) (int, int, bool) {
	parts := strings.Split(owner, ":")
	if len(parts) != 2 {
		return 0, 0, false
	}
	ids := make([]int, len(parts))
	for i, s := range parts {
		var err error
		ids[i], err = strconv.Atoi(s)
		if err != nil {
			return 0, 0, false
		}
	}
	return ids[0], ids[1], true
}

func (t *DeviceConfigurer) ConfigureDevices(instance *Instance) error {
	owner, err := instance.GetOwner()
	if err != nil {
		return err
	}
	source := instance.DeviceSource()
	if source.IsDefined() && source.Clone {
		err = t.CreateFilesystems(instance, source.Instance, source.Snapshot, owner)
	} else {
		err = t.CreateFilesystems(instance, nil, "", owner)
	}
	if err != nil {
		return err
	}
	filesystems, err := instance.Filesystems()
	if err != nil {
		return err
	}

	script := t.NewScript()
	devices := SortDevices(t.Config.Devices)
	for key, d := range devices {
		if d.Device.Filesystem == "" {
			continue
		}
		dir, err := instance.DeviceDir(d.Name, d.Device)
		if err != nil {
			return err
		}
		fs, found := filesystems[d.Device.Filesystem]
		if !found {
			return fmt.Errorf("missing filesystem: %s device: \n", d.Device.Filesystem, key)
		}
		if !fs.IsNew && util.DirExists(dir) {
			continue
		}
		err = t.CreateDir(dir, owner)
		if err != nil {
			return err
		}
		if !t.NoRsync && !fs.Filesystem.Transient {
			if source.IsDefined() && !source.Clone {
				templateDir, err := source.Instance.DeviceDir(d.Name, d.Device)
				if err != nil {
					return err
				}
				if templateDir != "" && util.DirExists(templateDir) {
					script.Run("sudo", "rsync", "-a", templateDir+"/", dir+"/")
				} else {
					fmt.Printf("skipping missing template Device=%s dir=%s\n", d.Name, templateDir)
				}
			}
		}
		if script.Error() != nil {
			return script.Error()
		}
	}
	return nil
}

func (t *DeviceConfigurer) ExtractDevices(instance *Instance, server srv.InstanceServer) error {
	owner, err := instance.GetOwner()
	if err != nil {
		return err
	}
	err = t.CreateFilesystems(instance, nil, "", owner)
	if err != nil {
		return err
	}
	filesystems, err := instance.Filesystems()
	if err != nil {
		return err
	}

	script := t.NewScript()
	devices := SortDevices(t.Config.Devices)
	var uid string
	var gid string
	rootFS := NewRootFS(server, instance.Name)
	defer rootFS.Unmount()
	u, g, ok := parseOwner(owner)
	if !ok {
		return fmt.Errorf("owner should have the form uid:gid (%s)", owner)
	}
	uid = strconv.Itoa(u)
	gid = strconv.Itoa(g)
	lxops, err := os.Executable()
	if err != nil {
		return fmt.Errorf("cannot locate executable")
	}
	for key, d := range devices {
		if d.Device.Filesystem == "" {
			continue
		}
		dir, err := instance.DeviceDir(d.Name, d.Device)
		if err != nil {
			return err
		}
		fs, found := filesystems[d.Device.Filesystem]
		if !found {
			return fmt.Errorf("missing filesystem: %s device: \n", d.Device.Filesystem, key)
		}
		if !fs.IsNew && util.DirExists(dir) {
			continue
		}
		if !fs.Filesystem.Transient {
			fmt.Printf("using %s from image\n", d.Device.Path)
			err := rootFS.Mount()
			if err != nil {
				return err
			}
			mountedDir, err := rootFS.MountedDir(d.Device.Path)
			if err != nil {
				return err
			}
			script.Run("sudo", "rsync", "-a", mountedDir+"/", dir+"/")
			script.Run("sudo", lxops, "shiftids", "-u", uid, "-g", gid, "-v", dir)
		}
		if script.Error() != nil {
			return script.Error()
		}
	}
	return nil
}

func (t *DeviceConfigurer) CreateProfile(client srv.Client, instance *Instance) error {
	profileName := instance.ProfileName()
	if profileName == "" {
		return nil
	}
	devices, err := instance.NewDeviceMap()
	if err != nil {
		return err
	}

	server, err := client.ProjectInstanceServer(t.Config.Project)
	if err != nil {
		return err
	}
	profile := &srv.Profile{
		Name:        profileName,
		Description: t.Config.AbsFilename(),
		Devices:     devices,
	}

	if t.Trace {
		instance.PrintDevices()
	}
	if !t.DryRun {
		return server.CreateProfile(profile)
	}
	return nil
}

func (t *DeviceConfigurer) RenameFilesystems(oldInstance, newInstance *Instance) error {
	oldPaths, err := oldInstance.FilesystemList()
	if err != nil {
		return err
	}
	newPaths, err := newInstance.Filesystems()
	if err != nil {
		return err
	}
	s := t.NewScript()
	for _, oldpath := range InstanceFSList(oldPaths).Roots() {
		newpath := newPaths[oldpath.Id]
		if oldpath.Path == newpath.Path {
			continue
		}
		if oldpath.IsDir() {
			newdir := newpath.Dir()
			if util.DirExists(newdir) {
				return errors.New(newdir + ": already exists")
			}
			s.Run("mv", oldpath.Dir(), newdir)
		} else {
			s.Run("sudo", "zfs", "rename", oldpath.Path, newpath.Path)
		}
	}
	return s.Error()
}
