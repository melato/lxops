module melato.org/lxops

go 1.18

replace (
	melato.org/cloudconfig => ../cloudconfig
	melato.org/cloudconfiglxd => ../cloudconfiglxd
)

require (
	gopkg.in/yaml.v2 v2.4.0
	melato.org/cloudconfig v0.0.0-00010101000000-000000000000
	melato.org/command v1.0.1
	melato.org/script v1.0.0
	melato.org/table3 v0.0.0-20220501091508-83fb75c200b0
)
