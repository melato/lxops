package lxops

import (
	"melato.org/lxops/cfg"
	"melato.org/lxops/cfg/migrate"
)

func InitConfigTypes() {
	cfg.SetMigrateFunc("#new", migrate.MigrateNew)
}
