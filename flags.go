package lxops

type DryRunFlag struct {
	DryRun bool `name:"-" usage:"show the commands to run, but do not change anything"`
}
