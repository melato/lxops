package lxops

import (
	"melato.org/cloudconfig/ostype"
)

func InitOSTypes() {
	OSTypes["alpine"] = &ostype.Alpine{}
	OSTypes["debian"] = &ostype.Debian{}
	OSTypes["ubuntu"] = &ostype.Debian{}
}
