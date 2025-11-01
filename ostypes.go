package lxops

import (
	"melato.org/cloudconfig/ostype"
	"melato.org/lxops/cfg"
	lxostype "melato.org/lxops/ostype"
)

func InitOSTypes() {
	cfg.OSTypes["alpine"] = &ostype.Alpine{}
	cfg.OSTypes["debian"] = &ostype.Debian{}
	cfg.OSTypes["ubuntu"] = &ostype.Debian{}
	cfg.OSTypes["openwrt"] = &lxostype.Openwrt{}
}
