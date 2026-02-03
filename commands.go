package lxops

import (
	"fmt"
	"io/fs"
	"os"

	"melato.org/command"
	"melato.org/command/usage"
	"melato.org/lxops/cfg"
	"melato.org/lxops/cli"
	"melato.org/lxops/help"
	"melato.org/lxops/internal/templatefs"
	"melato.org/lxops/srv"
)

func helpDataModel(serverType string) any {
	return map[string]string{
		"ServerType":    serverType,
		"ConfigVersion": cfg.Comment,
	}
}

func helpFS(serverType string) fs.FS {
	var fsys fs.FS = help.FS
	helpDir, ok := os.LookupEnv("LXOPS_HELP")
	if ok {
		fsys = os.DirFS(helpDir)
	}
	return templatefs.NewTemplateFS(fsys, helpDataModel(serverType))
}

func RootCommand(client srv.Client) *command.SimpleCommand {
	serverType := client.ServerType()
	helpFS := helpFS(serverType)
	usageData, err := fs.ReadFile(helpFS, "commands.yaml.tpl")
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	var cmd command.SimpleCommand
	cmd.Flags(client)
	launcher := &Launcher{Client: client}
	//cmd.Command("create").Flags(launcher).RunFunc(launcher.InstanceFunc(launcher.CreateContainer, true))
	cmd.Command("extract").Flags(launcher).RunFunc(launcher.InstanceFunc(launcher.ExtractDevices, true))
	cmd.Command("launch").Flags(launcher).RunFunc(launcher.InstanceFunc(launcher.LaunchContainer, true))
	cmd.Command("delete").Flags(launcher).RunFunc(launcher.InstanceFunc(launcher.DeleteContainer, false))
	cmd.Command("destroy").Flags(launcher).RunFunc(launcher.InstanceFunc(launcher.DestroyContainer, false))
	cmd.Command("rebuild").Flags(launcher).RunFunc(launcher.InstanceFunc(launcher.Rebuild, true))
	cmd.Command("rename").Flags(launcher).RunFunc(launcher.Rename)
	cmd.Command("create-devices").Flags(launcher).RunFunc(launcher.InstanceFunc(launcher.CreateDevices, true))
	cmd.Command("create-profile").Flags(launcher).RunFunc(launcher.InstanceFunc(launcher.CreateProfile, false))

	cmd.Command("ostypes").RunFunc(cli.ListOSTypes)

	snapshot := &Snapshot{Client: client}
	cmd.Command("snapshot").Flags(snapshot).RunFunc(snapshot.InstanceFunc(snapshot.Run, false))

	rollback := &Rollback{Client: client}
	cmd.Command("rollback").Flags(rollback).RunFunc(rollback.InstanceFunc(rollback.Run, false))

	configurer := &Configurer{Client: client}
	cmd.Command("configure").Flags(configurer).RunFunc(configurer.InstanceFunc(configurer.ConfigureContainer, false))

	instanceOps := &InstanceOps{Client: client}
	instanceCmd := cmd.Command("i").Flags(instanceOps)
	instanceCmd.Command("verify").RunFunc(instanceOps.InstanceFunc(instanceOps.Verify, true))
	instanceCmd.Command("description").RunFunc(instanceOps.InstanceFunc(instanceOps.Description, false))
	instanceCmd.Command("properties").RunFunc(instanceOps.InstanceFunc(instanceOps.Properties, false))
	instanceCmd.Command("filesystems").RunFunc(instanceOps.InstanceFunc(instanceOps.Filesystems, false))
	instanceCmd.Command("devices").RunFunc(instanceOps.InstanceFunc(instanceOps.Devices, false))
	instanceCmd.Command("project").RunFunc(instanceOps.InstanceFunc(instanceOps.Project, false))

	profile := cmd.Command("profile")
	profileConfigurer := &ProfileConfigurer{Client: client}
	profile.Command("list").Flags(profileConfigurer).RunFunc(profileConfigurer.InstanceFunc(profileConfigurer.List, false))
	profile.Command("diff").Flags(profileConfigurer).RunFunc(profileConfigurer.InstanceFunc(profileConfigurer.Diff, false))
	profile.Command("apply").Flags(profileConfigurer).RunFunc(profileConfigurer.InstanceFunc(profileConfigurer.Apply, false))
	profile.Command("reorder").Flags(profileConfigurer).RunFunc(profileConfigurer.InstanceFunc(profileConfigurer.Reorder, false))
	profileOps := &cli.ProfileOps{Client: client}
	profile.Command("import").Flags(profileOps).RunFunc(profileOps.Import)
	profile.Command("exists").RunFunc(profileOps.ProfileExists)
	profileExportOps := &cli.ProfileExportOps{Client: client}
	profile.Command("export").Flags(profileExportOps).RunFunc(profileExportOps.Export)

	propertyOps := &PropertyOptions{}
	propertyCmd := cmd.Command("property").Flags(propertyOps)
	propertyCmd.Command("list").RunFunc(propertyOps.List)
	propertyCmd.Command("set").RunFunc(propertyOps.Set)
	propertyCmd.Command("get").RunFunc(propertyOps.Get)
	propertyCmd.Command("file").RunFunc(propertyOps.File)

	parse := &ParseOps{}
	configCmd := cmd.Command("config").Flags(parse)
	configCmd.Command("parse").RunFunc(parse.Parse)
	configCmd.Command("print").RunFunc(parse.Print)
	configOps := &ConfigOps{}
	configCmd.Command("formats").RunFunc(configOps.Formats)
	configCmd.Command("properties").RunFunc(parse.PrintProperties)
	configCmd.Command("packages").RunFunc(parse.PrintPackages)
	configCmd.Command("cloud-config").RunFunc(parse.CloudConfigFiles)
	configCmd.Command("includes").RunFunc(configOps.Includes)
	convert := &cli.ConvertOps{}
	configCmd.Command("convert").Flags(convert).RunFunc(convert.Convert)

	containerOps := &cli.InstanceOps{Client: client}
	containerCmd := cmd.Command("instance")
	containerCmd.Flags(containerOps)
	containerCmd.Command("profiles").RunFunc(containerOps.Profiles)
	containerCmd.Command("wait").RunFunc(containerOps.Wait)
	containerCmd.Command("devices").RunFunc(containerOps.Devices)
	containerCmd.Command("hwaddr").RunFunc(containerOps.ListHwaddr)
	containerCmd.Command("info").RunFunc(containerOps.Info)

	cloudconfig := &cli.Cloudconfig{Client: client}
	cloudconfigInstanceOps := &cli.CloudconfigInstanceOps{Cloudconfig: cloudconfig}
	containerCmd.Command("cloudconfig").Flags(cloudconfigInstanceOps).RunFunc(cloudconfigInstanceOps.Apply)

	cloudconfigOps := &cli.CloudconfigOps{Cloudconfig: cloudconfig}
	cmd.Command("cloudconfig").Flags(cloudconfigOps).RunFunc(cloudconfigOps.Apply)

	imageCmd := cmd.Command("image")
	imageCmd.Command("instances").Flags(containerOps).RunFunc(containerOps.ListImages)

	publishOps := &PublishOps{InstanceOps: cli.InstanceOps{Client: client}}
	cmd.Command("publish").Flags(publishOps).RunFunc(publishOps.PublishInstance)

	networkOp := &cli.NetworkOp{Client: client}
	containerCmd.Command("addresses").Flags(networkOp).RunFunc(networkOp.ExportAddresses)

	numberOp := &cli.AssignNumbers{Client: client}
	containerCmd.Command("number").Flags(numberOp).RunFunc(numberOp.Run)

	exportOps := &ExportOps{Client: client}
	cmd.Command("export").Flags(exportOps).RunFunc(exportOps.Export)
	cmd.Command("import").Flags(exportOps).RunFunc(exportOps.Import)

	var shiftIds cli.ShiftIds
	cmd.Command("shiftids").Flags(&shiftIds).RunFunc(shiftIds.Run)

	var migrate Migrate
	cmd.Command("copy-filesystems").Flags(&migrate).RunFunc(migrate.CopyFilesystems)

	usage.Apply(&cmd, usageData)

	return &cmd
}
