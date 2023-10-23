package lxops

import (
	"melato.org/cloudconfig/ostype"
	"melato.org/lxops/cfg"
)

func InitOSTypes() {
	cfg.OSTypes["alpine"] = &ostype.Alpine{}
	cfg.OSTypes["debian"] = &ostype.Debian{}
	cfg.OSTypes["ubuntu"] = &ostype.Debian{}
}
