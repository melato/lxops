package lxdutil

import (
	"fmt"
	"strings"

	"github.com/canonical/lxd/shared/api"
	"melato.org/lxops/srv"
	"melato.org/lxops/yaml"
	"melato.org/script"
)

var Trace bool = true

func (t *InstanceServer) GetProfileNames() ([]string, error) {
	return t.Server.GetProfileNames()
}

/*
func (t *InstanceServer) ProfileExists(name string) (bool, error) {
	profiles, err := t.GetProfileNames()
	if err != nil {
		return false, err
	}
	for _, profile := range profiles {
		if profile == name {
			return true, nil
		}
	}
	return false, nil
}
*/

func DevicesToMap(devices map[string]*srv.Device) map[string]map[string]string {
	m := make(map[string]map[string]string)

	for deviceName, device := range devices {

		d := map[string]string{
			"type":   "disk",
			"path":   device.Path,
			"source": device.Source,
		}
		if device.Readonly {
			d["readonly"] = "true"
		}
		m[deviceName] = d
	}
	return m

}
func (t *InstanceServer) CreateProfile(profile *srv.Profile) error {
	devices := DevicesToMap(profile.Devices)

	post := api.ProfilesPost{Name: profile.Name, ProfilePut: api.ProfilePut{
		Devices:     devices,
		Config:      profile.Config,
		Description: profile.Description}}
	return t.Server.CreateProfile(post)
}

func (t *InstanceServer) DeleteProfile(profile string) error {
	err := t.Server.DeleteProfile(profile)
	if err != nil {
		return fmt.Errorf("delete profile %s: %w", profile, err)
	}
	return nil
}

func (t *InstanceServer) LaunchInstance(launch *srv.Launch) error {
	var lxcArgs []string
	if launch.Project != "" {
		lxcArgs = append(lxcArgs, "--project", launch.Project)
	}
	lxcArgs = append(lxcArgs, "init")

	lxcArgs = append(lxcArgs, launch.Image)
	for _, profile := range launch.Profiles {
		lxcArgs = append(lxcArgs, "-p", profile)
	}
	for _, option := range launch.LxcOptions {
		lxcArgs = append(lxcArgs, option)
	}
	lxcArgs = append(lxcArgs, launch.Name)
	s := &script.Script{Trace: Trace}

	s.Run("lxc", lxcArgs...)
	if s.HasError() {
		return s.Error()
	}
	err := t.SetInstanceNetwork(launch.Name, launch.Network)
	if err != nil {
		return err
	}
	err = t.StartInstance(launch.Name)
	if err != nil {
		return err
	}

	err = WaitForNetwork(t.Server, launch.Name)
	if err != nil {
		return err
	}
	return nil

}

func (t *InstanceServer) RebuildInstance(image, instance string) error {
	s := &script.Script{Trace: Trace}

	s.Run("lxc", "rebuild", "-f", image, instance)
	if s.HasError() {
		return s.Error()
	}

	state, _, err := t.Server.GetInstanceState(instance)
	if err != nil {
		return err
	}

	if state.Status != Running {
		err = t.StartInstance(instance)
		if err != nil {
			return err
		}
	}

	err = WaitForNetwork(t.Server, instance)
	if err != nil {
		return err
	}
	return nil
}

func (t *InstanceServer) CopyInstance(cp *srv.Copy) error {
	var copyArgs []string
	if cp.Project != "" {
		copyArgs = append(copyArgs, "--project", cp.Project)
	}

	copyArgs = append(copyArgs, "copy")

	if cp.Project != "" {
		copyArgs = append(copyArgs, "--target-project", cp.Project)
	}
	if cp.SourceSnapshot == "" {
		copyArgs = append(copyArgs, "--instance-only", cp.SourceInstance)
	} else {
		copyArgs = append(copyArgs, cp.SourceInstance+"/"+cp.SourceSnapshot)
	}
	copyArgs = append(copyArgs, cp.Name)
	s := &script.Script{Trace: Trace}
	s.Run("lxc", copyArgs...)
	if s.HasError() {
		return s.Error()
	}
	err := t.SetInstanceNetwork(cp.Name, cp.Network)
	if err != nil {
		return err
	}
	err = t.StartInstance(cp.Name)
	if err != nil {
		return err
	}

	err = WaitForNetwork(t.Server, cp.Name)
	if err != nil {
		return err
	}
	return nil
}

func (t *InstanceServer) NewConfigurer(instance string) (srv.InstanceConfigurer, error) {
	return NewInstanceConfigurer(t.Server, instance), nil
}

func (t *InstanceServer) GetInstanceProfiles(name string) ([]string, error) {
	c, _, err := t.Server.GetInstance(name)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", name, err)
	}
	return c.Profiles, nil
}

type Network struct {
	Hwaddresses map[string]string
}

func (t *InstanceServer) GetInstanceNetwork(name string) (srv.Network, error) {
	state, _, err := t.Server.GetContainerState(name)
	if err != nil {
		// assume container doesn't exist.  ignore error, empty network
		return nil, nil
	}
	var hwaddresses map[string]string
	for network, networkState := range state.Network {
		if networkState.Hwaddr == "" {
			continue
		}
		if hwaddresses == nil {
			hwaddresses = make(map[string]string)
		}
		hwaddresses[network] = networkState.Hwaddr
	}
	return &Network{hwaddresses}, nil
}

func (t *InstanceServer) SetInstanceNetwork(name string, v srv.Network) error {
	if v == nil {
		return nil
	}
	net, ok := v.(*Network)
	if !ok {
		return fmt.Errorf("invalid network: %T", v)
	}
	if len(net.Hwaddresses) == 0 {
		return nil
	}
	c, etag, err := t.Server.GetInstance(name)
	if err != nil {
		return fmt.Errorf("%s: %w", name, err)
	}
	for network, hwaddr := range net.Hwaddresses {
		key := "volatile." + network + ".hwaddr"
		c.InstancePut.Config[key] = hwaddr
		if Trace {
			fmt.Printf("set config %s: %s\n", key, hwaddr)
		}
	}
	op, err := t.Server.UpdateInstance(name, c.InstancePut, etag)
	if err != nil {
		return err
	}
	if err := op.Wait(); err != nil {
		return fmt.Errorf("%s: %w", name, err)
	}
	return nil
}

func (t *InstanceServer) WaitForNetwork(instance string) error {
	return WaitForNetwork(t.Server, instance)
}

// configureContainer configures the container directly, if necessary, and starts it
func (t *InstanceServer) configureContainer(launch *srv.Launch) error {
	return fmt.Errorf("not implemented")
}

func (t *InstanceServer) DeleteInstance(name string, stop bool) error {
	if stop {
		_ = t.StopInstance(name)
	}

	op, err := t.Server.DeleteInstance(name)
	if err == nil {
		if Trace {
			fmt.Printf("deleted instance %s\n", name)
		}
		if err := op.Wait(); err != nil {
			return fmt.Errorf("%s: %w", name, err)
		}
	} else {
		state, _, err := t.Server.GetInstanceState(name)
		if err == nil {
			return fmt.Errorf("instance %s is %s", name, state.Status)
		}
	}
	return nil
}

func (t *InstanceServer) SetInstanceProfiles(instance string, profiles []string) error {
	c, etag, err := t.Server.GetInstance(instance)
	if err != nil {
		return fmt.Errorf("%s: %w", instance, err)
	}
	c.Profiles = profiles
	op, err := t.Server.UpdateInstance(instance, c.InstancePut, etag)
	if err != nil {
		return err
	}
	if err := op.Wait(); err != nil {
		return fmt.Errorf("%s: %w", instance, err)
	}
	return nil
}

func (t *InstanceServer) CreateInstanceSnapshot(instance string, snapshot string) error {
	op, err := t.Server.CreateContainerSnapshot(instance, api.ContainerSnapshotsPost{Name: snapshot})
	if err != nil {
		return fmt.Errorf("%s: %w", instance, err)
	}
	if err := op.Wait(); err != nil {
		return fmt.Errorf("%s: %w", instance, err)
	}
	return nil
}

func (t *InstanceServer) RenameInstance(oldname, newname string) error {
	op, err := t.Server.RenameInstance(oldname, api.InstancePost{Name: newname})
	if err != nil {
		return fmt.Errorf("%s: %w", oldname, err)
	}
	if err := op.Wait(); err != nil {
		return err
	}
	return nil
}

func (t *InstanceServer) ExportProfile(name string) ([]byte, error) {
	profile, _, err := t.Server.GetProfile(name)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", err, name)
	}

	return yaml.Marshal(&profile.ProfilePut)
}

func (t *InstanceServer) ImportProfile(name string, data []byte) error {
	var profile api.ProfilePut
	err := yaml.Unmarshal(data, &profile)
	if err != nil {
		return err
	}
	exists, _ := t.ProfileExists(name)
	if exists {
		err = t.Server.UpdateProfile(name, profile, "")
	} else {
		var post api.ProfilesPost
		post.ProfilePut = profile
		post.Name = name
		err = t.Server.CreateProfile(post)
	}
	if err != nil {
		return fmt.Errorf("%s: %w", name, err)
	}
	return nil
}

func (t *InstanceServer) GetInstanceDevices(instance string) (map[string]*srv.Device, error) {
	c, _, err := t.Server.GetInstance(instance)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", instance, err)
	}
	devices := make(map[string]*srv.Device)
	for name, d := range c.ExpandedDevices {
		if d["type"] == "disk" {
			device := &srv.Device{Path: d["path"], Source: d["source"], Readonly: d["readonly"] == "true"}
			devices[name] = device
		}
	}
	return devices, nil
}

func (t *InstanceServer) GetHwaddresses() ([]srv.Hwaddr, error) {
	instances, err := t.Server.GetInstances(api.InstanceTypeAny)
	if err != nil {
		return nil, err
	}
	addresses := make([]srv.Hwaddr, 0, len(instances))
	var i api.Instance
	for _, i = range instances {
		addresses = append(addresses, srv.Hwaddr{Instance: i.Name, Hwaddr: i.Config["volatile.eth0.hwaddr"]})
	}
	return addresses, nil

}

func (t *InstanceServer) GetInstanceImages() ([]srv.InstanceImage, error) {
	images, err := t.Server.GetImages()
	if err != nil {
		return nil, err
	}
	fingerprints := make(map[string]string)
	for _, image := range images {
		names := make([]string, len(image.Aliases))
		for i, a := range image.Aliases {
			names[i] = a.Name
		}
		fingerprints[image.Fingerprint] = strings.Join(names, ",")
	}
	instances, err := t.Server.GetInstances(api.InstanceTypeAny)
	if err != nil {
		return nil, err
	}
	result := make([]srv.InstanceImage, 0, len(instances))
	for _, i := range instances {
		var im srv.InstanceImage
		im.Instance = i.Name
		fg := i.Config["volatile.base_image"]
		im.Image = fingerprints[fg]
		if im.Image == "" {
			im.Image = fg
		}
		result = append(result, im)
	}
	return result, nil

}

func (t *InstanceServer) GetInstanceNames(onlyRunning bool) ([]string, error) {
	containers, err := t.Server.GetContainersFull()
	if err != nil {
		return nil, err
	}
	var names []string
	for _, container := range containers {
		if onlyRunning && container.State.Status != Running {
			continue
		}
		names = append(names, container.Name)
	}
	return names, nil
}

func (t *InstanceServer) GetInstanceAddresses(family string) ([]*srv.HostAddress, error) {
	var addresses []*srv.HostAddress

	for _, instanceType := range []api.InstanceType{api.InstanceTypeContainer, api.InstanceTypeVM} {
		containers, err := t.Server.GetInstancesFull(instanceType)
		if err != nil {
			return nil, err
		}
		for _, c := range containers {
			if c.State == nil || c.State.Network == nil {
				continue
			}
			for _, net := range c.State.Network {
				for _, a := range net.Addresses {
					if a.Family == family && a.Scope == "global" {
						addresses = append(addresses, &srv.HostAddress{Name: c.Name, Address: a.Address})
					}
				}
			}
		}
	}
	return addresses, nil
}

func (t *InstanceServer) PublishInstance(instance, snapshot, alias string) error {
	s := &script.Script{Trace: Trace}
	s.Run("lxc", "publish", instance+"/"+snapshot, "--alias="+alias)
	return s.Error()
}

func (t *InstanceServer) ExportImage(image string, path string) error {
	s := &script.Script{Trace: Trace}
	s.Run("lxc", "image", "export", image, path)
	return s.Error()
}

func (t *InstanceServer) ImportImage(image string, path string) error {
	s := &script.Script{Trace: Trace}
	s.Run("lxc", "image", "import", path, "--alias="+image)
	return s.Error()
}
