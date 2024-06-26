package lxops

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"

	"melato.org/lxops/cfg"
	"melato.org/lxops/srv"
	"melato.org/lxops/util"
	"melato.org/script"
)

type Launcher struct {
	Client srv.Client `name:"-"`
	ConfigOptions
	WaitBeforeConfigure int  `name:"wait-configure" usage:"# seconds to wait before configuration"`
	PollConfigure       int  `name:"poll" usage:"max # seconds to poll before configuration"`
	WaitBeforeStop      int  `name:"wait-stop" usage:"# seconds to wait before stop or snapshot"`
	Trace               bool `name:"t" usage:"trace print what is happening"`
	DryRunFlag
}

func (t *Launcher) Init() error {
	t.WaitBeforeStop = 5
	t.PollConfigure = 20
	return t.ConfigOptions.Init()
}

func (t *Launcher) Configured() error {
	if t.DryRun {
		t.Trace = true
	}
	t.ConfigOptions.ConfigureProject(t.Client)
	return t.ConfigOptions.Configured()
}

func (t *Launcher) NewScript() *script.Script {
	return &script.Script{Trace: t.Trace, DryRun: t.DryRun}
}

func (t *Launcher) Rebuild(instance *Instance) error {
	t.Trace = true
	config := instance.Config
	source := instance.ContainerSource()
	if source.IsDefined() {
		return fmt.Errorf("cannot rebuild from origin")
	}
	image, err := config.Image.Substitute(instance.Properties)
	if err != nil {
		return err
	}
	if image == "" {
		return fmt.Errorf("cannot rebuild without image")
	}
	server, err := t.Client.ProjectInstanceServer(config.Project)
	if err != nil {
		return err
	}
	err = server.RebuildInstance(image, instance.Name)
	if err != nil {
		return err
	}
	configurer := t.NewConfigurer()
	return configurer.ConfigureContainer(instance)
}

func (t *Launcher) NewConfigurer() *Configurer {
	var c = &Configurer{Client: t.Client, Trace: t.Trace, DryRunFlag: t.DryRunFlag}
	c.PollSeconds = t.PollConfigure
	return c
}

func (t *Launcher) lxcLaunch(instance *Instance, server srv.InstanceServer, options *launch_options) error {
	config := instance.Config
	_, err := cfg.OSType(config.Ostype)
	if err != nil {
		return err
	}
	var launch srv.Launch
	launch.Name = instance.Container()
	image, err := config.Image.Substitute(instance.Properties)
	if err != nil {
		return err
	}

	if image == "" {
		return errors.New("Please provide image or version")
	}
	launch.Image = image
	launch.Profiles = options.Profiles
	launch.LxcOptions = config.LxcOptions
	return server.LaunchInstance(&launch)
}

func (t *Launcher) createEmptyProfile(server srv.InstanceServer, profile string) error {
	return server.CreateProfile(&srv.Profile{
		Name:        profile,
		Description: "lxops placeholder profile",
	})
}

func (t *Launcher) deleteProfiles(server srv.InstanceServer, profiles []string) error {
	// delete the missing profiles from the new container, and delete them
	for _, profile := range profiles {
		if t.Trace {
			fmt.Printf("delete profile %s\n", profile)
		}
		if !t.DryRun {
			err := server.DeleteProfile(profile)
			if err != nil {
				return nil
			}
		}
	}
	return nil
}

type launch_options struct {
	Profiles []string
}

func (t *Launcher) copyContainer(instance *Instance, source ContainerSource, server srv.InstanceServer, options *launch_options) error {
	container := instance.Container()
	sourceServer, err := t.Client.ProjectInstanceServer(source.Project)
	if err != nil {
		return err
	}
	allProfiles, err := server.GetProfileNames()
	if err != nil {
		return err
	}
	instanceProfiles, err := sourceServer.GetInstanceProfiles(source.Container)
	if err != nil {
		return fmt.Errorf("%s_%s: %v", source.Project, source.Container, err)
	}
	missingProfiles := util.StringSlice(instanceProfiles).Diff(allProfiles)
	// lxc copy will fail if the source container has profiles that do not exist in the target server
	// so create the missing profiles, and delete them after the copy
	for _, profile := range missingProfiles {
		err := t.createEmptyProfile(server, profile)
		if err != nil {
			return err
		}
	}

	var cp srv.Copy
	cp.Name = container
	cp.Project = source.Project
	cp.SourceInstance = source.Container
	cp.SourceSnapshot = source.Snapshot
	cp.Profiles = options.Profiles

	err = server.CopyInstance(&cp)

	err2 := t.deleteProfiles(server, missingProfiles)
	if err != nil {
		return err
	}
	return err2
}

func (t *Launcher) CreateDevices(instance *Instance) error {
	t.Trace = true
	dev := NewDeviceConfigurer(instance.Config)
	dev.Trace = t.Trace
	dev.DryRun = t.DryRun
	return dev.ConfigureDevices(instance)
}

func (t *Launcher) CreateProfile(instance *Instance) error {
	dev := NewDeviceConfigurer(instance.Config)
	dev.Trace = t.Trace
	dev.DryRun = t.DryRun
	profileName := instance.ProfileName()
	if profileName != "" {
		fmt.Println(profileName)
		return dev.CreateProfile(t.Client, instance)
	} else {
		fmt.Printf("skipping instance %s: no lxops profile\n", instance.Name)
		return nil
	}

}

func (t *Launcher) LaunchContainer(instance *Instance) error {
	return t.launchContainer(instance)
}

func (t *Launcher) CreateContainer(instance *Instance) error {
	fmt.Println("create", instance.Name)
	t.Trace = true

	source := instance.ContainerSource()
	fmt.Printf("source:%v\n", source)
	if source.IsDefined() {
		return fmt.Errorf("cannot create container without image")
	}

	var err error

	config := instance.Config
	var create srv.Create
	create.Name = instance.Container()
	create.Image, err = config.Image.Substitute(instance.Properties)
	if err != nil {
		return err
	}
	if create.Image == "" {
		return errors.New("Please provide image")
	}
	create.LxcOptions = config.LxcOptions
	create.Profiles = append(create.Profiles, config.Profiles...)

	server, err := t.Client.ProjectInstanceServer(config.Project)
	if err != nil {
		return err
	}

	err = t.verifyProfiles(server, config.Profiles)
	if err != nil {
		return err
	}

	profileName := instance.ProfileName()
	if profileName == "" {
		// revisit this, if necessary
		return fmt.Errorf("configuration without profile is not supported")
	}

	err = server.CreateInstance(&create)

	if err != nil {
		return err
	}

	dev := NewDeviceConfigurer(instance.Config)
	dev.Trace = t.Trace
	dev.DryRun = t.DryRun
	err = dev.ConfigureDevices(instance)
	if err != nil {
		return err
	}

	if profileName != "" {
		err = dev.CreateProfile(t.Client, instance)
		if err != nil {
			return err
		}
		profiles := append(create.Profiles, profileName)
		err = server.SetInstanceProfiles(create.Name, profiles)
		if err != nil {
			return err
		}
	}

	return nil
}

func (t *Launcher) ExtractDevices(instance *Instance) error {
	fmt.Println("extract devices", instance.Name)
	t.Trace = true

	var err error

	config := instance.Config
	var create srv.Create
	create.Name = instance.Name
	create.Image, err = config.Image.Substitute(instance.Properties)
	if err != nil {
		return err
	}
	if create.Image == "" {
		return errors.New("Please provide image")
	}
	create.LxcOptions = config.LxcOptions
	create.Profiles = append(create.Profiles, config.Profiles...)

	server, err := t.Client.ProjectInstanceServer(config.Project)
	if err != nil {
		return err
	}

	err = t.verifyProfiles(server, config.Profiles)
	if err != nil {
		return err
	}

	err = server.CreateInstance(&create)

	if err != nil {
		return err
	}

	var deleted bool
	defer func() {
		if !deleted {
			server.DeleteInstance(instance.Name, false)
		}
	}()
	dev := NewDeviceConfigurer(instance.Config)
	dev.Trace = t.Trace
	dev.DryRun = t.DryRun
	err = dev.ExtractDevices(instance, server)
	if err != nil {
		return err
	}

	err = server.DeleteInstance(instance.Name, false)
	if err != nil {
		return err
	}
	deleted = true

	return nil
}

func (t *Launcher) verifyProfiles(server srv.InstanceServer, profiles []string) error {
	if len(profiles) == 0 {
		return nil
	}
	serverProfiles, err := server.GetProfileNames()
	if err != nil {
		return err
	}
	serverProfilesSet := make(util.Set[string])
	for _, profile := range serverProfiles {
		serverProfilesSet.Put(profile)
	}
	var missing []string
	for _, profile := range profiles {
		if !serverProfilesSet.Contains(profile) {
			missing = append(missing, profile)
		}
	}
	if len(missing) > 0 {
		return fmt.Errorf("missing profiles: %v", missing)
	}
	return nil
}

func (t *Launcher) launchContainer(instance *Instance) error {
	fmt.Println("launch", instance.Name)
	t.Trace = true
	config := instance.Config
	server, err := t.Client.ProjectInstanceServer(config.Project)
	if err != nil {
		return err
	}

	err = t.verifyProfiles(server, config.Profiles)
	if err != nil {
		return err
	}

	dev := NewDeviceConfigurer(instance.Config)
	dev.Trace = t.Trace
	dev.DryRun = t.DryRun
	err = dev.ConfigureDevices(instance)
	if err != nil {
		return err
	}

	profileName := instance.ProfileName()
	if profileName == "" {
		// revisit this, if necessary
		return fmt.Errorf("configuration without profile is not supported")
	}
	if profileName != "" {
		err = dev.CreateProfile(t.Client, instance)
		if err != nil {
			return err
		}
	}

	var profiles []string
	profiles = append(profiles, config.Profiles...)
	if config.Devices != nil {
		if len(profiles) == 0 {
			profiles = append(profiles, "default")
		}
		if profileName != "" {
			profiles = append(profiles, profileName)
		}
	}
	options := &launch_options{Profiles: profiles}
	container := instance.Container()
	source := instance.ContainerSource()
	fmt.Printf("source:%v\n", source)
	if !source.IsDefined() {
		err := t.lxcLaunch(instance, server, options)
		if err != nil {
			return err
		}
	} else {
		err := t.copyContainer(instance, *source, server, options)
		if err != nil {
			return err
		}
	}
	configurer := t.NewConfigurer()
	if t.WaitBeforeConfigure != 0 {
		fmt.Printf("waiting %d seconds before configuring\n", t.WaitBeforeConfigure)
		time.Sleep(time.Duration(t.WaitBeforeConfigure) * time.Second)
	}
	err = configurer.ConfigureContainer(instance)
	if err != nil {
		return err
	}
	if config.Stop || config.Snapshot != "" {
		if t.WaitBeforeStop != 0 {
			fmt.Printf("waiting %d seconds for container installation scripts to complete\n", t.WaitBeforeStop)
			time.Sleep(time.Duration(t.WaitBeforeStop) * time.Second)
		}
	}
	if config.Stop {
		fmt.Printf("stop %s\n", container)
		if !t.DryRun {
			err = server.StopInstance(container)
			if err != nil {
				return err
			}
		}
	}
	if config.Snapshot != "" {
		fmt.Printf("snapshot %s %s\n", container, config.Snapshot)
		if !t.DryRun {
			err := server.CreateInstanceSnapshot(container, config.Snapshot)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (t *Launcher) deleteContainer(instance *Instance, stop bool) error {
	config := instance.Config
	container := instance.Container()
	server, err := t.Client.ProjectInstanceServer(config.Project)
	if err != nil {
		return err
	}
	if !t.DryRun {
		err = server.DeleteInstance(container, stop)
		if t.Trace {
			if err != nil {
				fmt.Printf("DeleteInstance(%s): %v\n", container, err)
			} else {
				fmt.Printf("deleted instance %s\n", container)
			}
		}
		if err != nil && t.Trace {
		}
	}
	profileName := instance.ProfileName()

	if !t.DryRun {
		err := server.DeleteProfile(profileName)
		if err != nil {
			fmt.Printf("DeleteProfile(%s): %v\n", profileName, err)
		} else {
			fmt.Printf("deleted profile %s\n", profileName)
		}
	}
	return nil
}

func (t *Launcher) DeleteContainer(instance *Instance) error {
	err := t.deleteContainer(instance, false)
	if err != nil {
		return err
	}
	filesystems, err := instance.FilesystemList()
	if err != nil {
		return err
	}
	var zfsFilesystems []string
	var dirFilesystems []string
	for _, fs := range filesystems {
		if fs.IsZfs() {
			zfsFilesystems = append(zfsFilesystems, fs.Path)
		} else {
			dirFilesystems = append(dirFilesystems, fs.Path)
		}
	}
	fmt.Fprintln(os.Stderr, "remaining filesystems:")
	if len(zfsFilesystems) > 0 {
		cmd := exec.Command("zfs", append([]string{"list", "-o", "name,used,referenced,origin,mountpoint"}, zfsFilesystems...)...)
		cmd.Stderr = io.Discard
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	if len(dirFilesystems) > 0 {
		cmd := exec.Command("ls", append([]string{"-l"}, dirFilesystems...)...)
		cmd.Stderr = io.Discard
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	return nil
}

func (t *Launcher) DestroyContainer(instance *Instance) error {
	err := t.deleteContainer(instance, false)
	if err != nil {
		return err
	}
	filesystems, err := instance.FilesystemList()
	if err != nil {
		return err
	}
	var zfsFilesystems []string
	var dirFilesystems []string
	for _, fs := range filesystems {
		if fs.Filesystem.Destroy {
			if fs.IsZfs() {
				zfsFilesystems = append(zfsFilesystems, fs.Path)
			} else {
				dirFilesystems = append(dirFilesystems, fs.Path)
			}
		}
	}
	if len(zfsFilesystems) > 0 {
		s := script.Script{DryRun: t.DryRun, Trace: t.Trace}
		lines := s.Cmd("zfs", append([]string{"list", "-H", "-o", "name"}, zfsFilesystems...)...).ToLines()
		s.Errors.Clear()
		for _, line := range lines {
			s.Run("sudo", "zfs", "destroy", "-r", line)
		}
		if s.HasError() {
			return s.Error()
		}
	}
	var firstError error
	for _, dir := range dirFilesystems {
		err := os.RemoveAll(dir)
		if err != nil && firstError == nil {
			firstError = err
		}
	}
	return firstError
}

func (t *Launcher) Rename(configFile string, newname string) error {
	instance, err := t.ConfigOptions.Instance(configFile)
	if err != nil {
		return err
	}

	if t.Trace {
		fmt.Printf("rename container %s -> %s\n", instance.Name, newname)
	}
	if instance.Name == newname {
		return errors.New("cannot rename to the same name")
	}
	oldprofile := instance.ProfileName()
	newInstance, err := instance.NewInstance(newname)
	if err != nil {
		return err
	}
	newprofile := newInstance.ProfileName()
	server, err := t.Client.ProjectInstanceServer(instance.Config.Project)
	if err != nil {
		return err
	}
	dev := NewDeviceConfigurer(instance.Config)
	dev.Trace = t.Trace
	dev.DryRun = t.DryRun

	containerName := instance.Container()
	newContainerName := newInstance.Container()
	var profiles []string
	if len(instance.Config.Devices) > 0 {
		profileExists, err := server.ProfileExists(newprofile)
		if profileExists {
			return errors.New(fmt.Sprintf("profile %s already exists", newprofile))
		}
		profiles, err = server.GetInstanceProfiles(containerName)
		if err != nil {
			return err
		}
	}
	if !t.DryRun {
		err := server.RenameInstance(containerName, newContainerName)
		if err != nil {
			return err
		}
	}
	if len(instance.Config.Devices) > 0 {
		err = dev.RenameFilesystems(instance, newInstance)
		if err != nil {
			return err
		}
		err = dev.CreateProfile(t.Client, newInstance)
		if err != nil {
			return err
		}
		var replaced bool
		for i, profile := range profiles {
			if profile == oldprofile {
				profiles[i] = newprofile
				replaced = true
				break
			}
		}
		if !replaced {
			profiles = append(profiles, newprofile)
		}
		if t.Trace {
			fmt.Printf("apply %s profiles: %v\n", newname, profiles)
		}
		if !t.DryRun {
			err := server.SetInstanceProfiles(newContainerName, profiles)
			if err != nil {
				return err
			}
		}
		if t.Trace {
			fmt.Printf("delete profile %s\n", oldprofile)
		}
		if !t.DryRun {
			err = server.DeleteProfile(oldprofile)
			if err != nil {
				return nil
			}
		}
	}
	return nil
}
