package ostype

import (
	"melato.org/cloudconfig"
)

type Openwrt struct {
}

func (t *Openwrt) NeedUserPasswords() bool { return false }

func (t *Openwrt) UpdatePackagesCommand() string { return "apk update" }

func (t *Openwrt) UpgradePackagesCommand() string { return "apk upgrade" }

func (t *Openwrt) InstallPackageCommand(pkg string) string {
	return "apk add " + pkg
}

func (t *Openwrt) AddUserCommand(u *cloudconfig.User) []string {
	return nil
}

func (t *Openwrt) SetTimezoneCommand(timezone string) []string {
	return nil
}
