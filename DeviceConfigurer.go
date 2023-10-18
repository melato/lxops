package lxops

import (
	"errors"
	"fmt"
	"os"

	"melato.org/lxops/srv"
	"melato.org/lxops/util"
	"melato.org/script"
)

type DeviceConfigurer struct {
	Config *Config
	Owner  string
	// NoRsync - do not rsync devices.  Use when importing.
	NoRsync bool
	Trace   bool
	DryRun  bool
}

func NewDeviceConfigurer(instance *Instance) (*DeviceConfigurer, error) {
	t := &DeviceConfigurer{Config: instance.Config}
	var err error
	t.Owner, err = instance.Config.DeviceOwner.Substitute(instance.Properties)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (t *DeviceConfigurer) NewScript() *script.Script {
	return &script.Script{Trace: t.Trace, DryRun: t.DryRun}
}

func (t *DeviceConfigurer) chownDir(scr *script.Script, dir string) {
	//scr.Run("sudo", "chown", "1000000:1000000", dir)
	if t.Owner != "" {
		scr.Run("sudo", "chown", t.Owner, dir)
	}
}

func (t *DeviceConfigurer) CreateDir(dir string, chown bool) error {
	if !util.DirExists(dir) {
		script := t.NewScript()
		script.Run("sudo", "mkdir", "-p", dir)
		//err = os.Mkdir(dir, 0755)
		if chown {
			t.chownDir(script, dir)
		}
		return script.Error()
	}
	return nil
}

func (t *DeviceConfigurer) CreateFilesystem(fs *InstanceFS, originDataset string) error {
	if fs.IsDir() {
		fs.IsNew = true
		return t.CreateDir(fs.Dir(), true)
	}

	var args []string
	if originDataset != "" {
		args = []string{"zfs", "clone", "-p"}
	} else {
		args = []string{"zfs", "create", "-p"}

	}

	// add properties
	for key, value := range fs.Filesystem.Zfsproperties {
		args = append(args, "-o", key+"="+value)
	}

	if originDataset != "" {
		args = append(args, originDataset)
	}
	args = append(args, fs.Path)
	s := t.NewScript()
	s.Run("sudo", args...)
	if originDataset == "" {
		t.chownDir(s, fs.Dir())
		fs.IsNew = true
	}
	return s.Error()
}

func (t *DeviceConfigurer) CreateFilesystems(instance, origin *Instance, snapshot string) error {
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
		if path.IsZfs() {
			originPath, exists := originPaths[path.Id]
			if exists {
				originDataset = originPath.Path + "@" + snapshot
			}
		}
		err := t.CreateFilesystem(path, originDataset)
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *DeviceConfigurer) ConfigureDevices(instance *Instance) error {
	source := instance.DeviceSource()
	var err error
	if source.IsDefined() && source.Clone {
		err = t.CreateFilesystems(instance, source.Instance, source.Snapshot)
	} else {
		err = t.CreateFilesystems(instance, nil, "")
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
	for _, d := range devices {
		dir, err := instance.DeviceDir(d.Name, d.Device)
		if err != nil {
			return err
		}
		fs, found := filesystems[d.Device.Filesystem]
		if !found {
			fmt.Fprintf(os.Stderr, "missing filesystem: %s\n", d.Device.Filesystem)
			continue
		}
		if !fs.IsNew && util.DirExists(dir) {
			continue
		}
		err = t.CreateDir(dir, true)
		if err != nil {
			return err
		}
		if !t.NoRsync && !fs.Filesystem.Transient && source.IsDefined() && !source.Clone {
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
		Description: "lxdops profile",
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
