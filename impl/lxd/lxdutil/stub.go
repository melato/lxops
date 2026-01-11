package lxdutil

import (
	"fmt"

	"melato.org/lxops/srv"
)

// These methods are not implemented for LXD yet.

func (t *InstanceServer) GetInstanceImageFields(name string) (*srv.ImageFields, error) {
	return nil, fmt.Errorf("GetInstanceImageFields: Unimplemented method.")
}

func (t *InstanceServer) GetInstance(name string) (any, error) {
	return nil, fmt.Errorf("GetInstance: Unimplemented method.")
}

func (t *InstanceServer) PublishInstance2(instance, snapshot,
	alias string,
	fields srv.ImageFields,
	options srv.PublishOptions) error {
	return fmt.Errorf("GetInstance: Unimplemented method.")
}
