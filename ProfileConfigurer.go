package lxops

import (
	"fmt"
	"strings"

	"melato.org/lxops/srv"
	"melato.org/lxops/util"
)

type ProfileConfigurer struct {
	Client srv.Client `name:"-"`
	ConfigOptions
	Trace bool
	DryRunFlag
}

func (t *ProfileConfigurer) Init() error {
	return t.ConfigOptions.Init()
}

func (t *ProfileConfigurer) Configured() error {
	if t.DryRun {
		t.Trace = true
	}
	t.ConfigOptions.ConfigureProject(t.Client)
	return t.ConfigOptions.Configured()
}

func (t *ProfileConfigurer) Profiles(instance *Instance) ([]string, error) {
	profile := instance.ProfileName()
	profiles := instance.Config.Profiles
	return append(profiles, profile), nil
}

func (t *ProfileConfigurer) Diff(instance *Instance) error {
	container := instance.Container()
	server, err := t.Client.ProjectInstanceServer(instance.Config.Project)
	if err != nil {
		return err
	}
	cProfiles, err := server.GetInstanceProfiles(container)
	if err != nil {
		return err
	}
	profiles, err := t.Profiles(instance)
	if err != nil {
		return err
	}
	if util.StringSlice(profiles).Equals(cProfiles) {
		return nil
	}
	onlyInConfig := util.StringSlice(profiles).Diff(cProfiles)
	onlyInContainer := util.StringSlice(cProfiles).Diff(profiles)
	sep := " "
	if len(onlyInConfig) > 0 {
		fmt.Printf("%s profiles only in config: %s\n", container, strings.Join(onlyInConfig, sep))
	}
	if len(onlyInContainer) > 0 {
		fmt.Printf("%s profiles only in container: %s\n", container, strings.Join(onlyInContainer, sep))
	}
	if len(onlyInConfig) == 0 && len(onlyInContainer) == 0 {
		fmt.Printf("%s profiles are in different order: %s\n", container, strings.Join(profiles, sep))
	}
	return nil
}

func (t *ProfileConfigurer) Reorder(instance *Instance) error {
	container := instance.Container()
	server, err := t.Client.ProjectInstanceServer(instance.Config.Project)
	if err != nil {
		return err
	}
	cProfiles, err := server.GetInstanceProfiles(container)
	if err != nil {
		return err
	}
	profiles, err := t.Profiles(instance)
	if err != nil {
		return err
	}
	if util.StringSlice(profiles).Equals(cProfiles) {
		return nil
	}

	sortedProfiles := util.StringSlice(profiles).Sorted()
	sortedContainer := util.StringSlice(cProfiles).Sorted()
	if util.StringSlice(sortedProfiles).Equals(sortedContainer) {
		return server.SetInstanceProfiles(container, profiles)
	}
	fmt.Println("profiles differ: " + container)
	return nil
}

func (t *ProfileConfigurer) Apply(instance *Instance) error {
	container := instance.Container()
	server, err := t.Client.ProjectInstanceServer(instance.Config.Project)
	if err != nil {
		return err
	}
	profiles, err := t.Profiles(instance)
	if err != nil {
		return err
	}
	return server.SetInstanceProfiles(container, profiles)
}

func (t *ProfileConfigurer) List(instance *Instance) error {
	profiles, err := t.Profiles(instance)
	if err != nil {
		return err
	}
	for _, profile := range profiles {
		fmt.Println(profile)
	}
	return nil
}
