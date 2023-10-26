package lxops

import (
	"melato.org/lxops/cfg"
	"melato.org/lxops/cfg/migrate"
)

func InitConfigTypes() {
	cfg.SetMigrateFunc("#lxdops", migrate.MigrateLxdops)
}
