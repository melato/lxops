package ostype

import (
	"melato.org/cloudconfig"
)

type OpenwrtOpkg struct {
}

func (t *OpenwrtOpkg) NeedUserPasswords() bool { return false }

func (t *OpenwrtOpkg) InstallPackageCommand(pkg string) string {
	return "opkg install " + pkg
}

func (t *OpenwrtOpkg) AddUserCommand(u *cloudconfig.User) []string {
	return nil
}

func (t *OpenwrtOpkg) SetTimezoneCommand(timezone string) []string {
	return nil
}
