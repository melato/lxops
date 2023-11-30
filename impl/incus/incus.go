package lxops_incus

import (
	"errors"
	"fmt"
	"strings"
	"time"

	incus "github.com/lxc/incus/client"
	"github.com/lxc/incus/shared/api"
)

const Running = "Running"

type InstanceServer struct {
	Server incus.InstanceServer
}

func SplitSnapshotName(name string) (container, snapshot string) {
	i := strings.Index(name, "/")
	if i >= 0 {
		return name[0:i], name[i+1:]
	} else {
		return name, ""
	}
}

func WaitForNetwork(server incus.InstanceServer, instance string) error {
	start := time.Now()
	var status string
	for i := 0; i < 300; i++ {
		state, _, err := server.GetInstanceState(instance)
		if err != nil {
			return fmt.Errorf("%s: %w", instance, err)
		}
		if state == nil {
			continue
		}
		for _, net := range state.Network {
			for _, a := range net.Addresses {
				if a.Family == "inet" && a.Scope == "global" {
					fmt.Println(a.Address)
					if i > 0 {
						fmt.Printf("time: %0.3fs\n", time.Now().Sub(start).Seconds())
					}
					return nil
				}
			}
		}
		if state.Status != status {
			status = state.Status
			fmt.Printf("status: %s time: %0.3fs\n", status, time.Now().Sub(start).Seconds())
		}

		time.Sleep(1 * time.Second)
	}
	return errors.New("could not get ip address for: " + instance)
}

func (t InstanceServer) updateInstanceState(container string, action string) error {
	op, err := t.Server.UpdateInstanceState(container, api.InstanceStatePut{Action: action}, "")
	if err != nil {
		return fmt.Errorf("%s: %w", container, err)
	}
	if err := op.Wait(); err != nil {
		return fmt.Errorf("%s: %w", container, err)
	}
	return nil
}

func (t InstanceServer) StartInstance(container string) error {
	return t.updateInstanceState(container, "start")
}

func (t InstanceServer) StopInstance(container string) error {
	return t.updateInstanceState(container, "stop")
}
func (t InstanceServer) ProfileExists(profile string) (bool, error) {
	_, _, err := t.Server.GetProfile(profile)
	if err == nil {
		return true, nil
	}
	return false, nil
}
