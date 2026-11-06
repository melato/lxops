module melato.org/lxops

go 1.23

replace melato.org/configfile => ../configfile

require (
	gopkg.in/yaml.v2 v2.4.0
	melato.org/cloudconfig v0.4.0
	melato.org/command v1.0.1
	melato.org/configfile v0.0.0-00010101000000-000000000000
	melato.org/script v1.0.0
	melato.org/table3 v0.1.1
)

require melato.org/trace v1.2.1
