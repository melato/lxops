package lxops

import (
	"bytes"
	_ "embed"
	"fmt"
	"os"
	"text/template"

	"melato.org/command"
	"melato.org/command/usage"
	"melato.org/lxops/cfg"
	"melato.org/lxops/cli"
	"melato.org/lxops/srv"
)

//go:embed commands.yaml.tpl
var usageTemplate string

func usageForServerType(serverType string) ([]byte, error) {
	data := usageTemplate
	envVar := "USAGE"
	if envVar != "" {
		file, ok := os.LookupEnv(envVar)
		if ok {
			if _, err := os.Stat(file); err == nil {
				fileContent, err := os.ReadFile(file)
				if err == nil {
					data = string(fileContent)
				} else {
					return nil, fmt.Errorf("%s: %v\n", file, err)
				}
			}
		}
	}
	tpl, err := template.New("").Parse(data)
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	err = tpl.Execute(&buf, map[string]string{
		"ServerType":    serverType,
		"ConfigVersion": cfg.Comment,
	})
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func RootCommand(client srv.Client) *command.SimpleCommand {
	usageData, err := usageForServerType(client.ServerType())
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	var cmd command.SimpleCommand
	cmd.Flags(client)
	launcher := &Launcher{Client: client}
	cmd.Command("launch").Flags(launcher).RunFunc(launcher.InstanceFunc(launcher.LaunchContainer, true))
	cmd.Command("delete").Flags(launcher).RunFunc(launcher.InstanceFunc(launcher.DeleteContainer, false))
	cmd.Command("destroy").Flags(launcher).RunFunc(launcher.InstanceFunc(launcher.DestroyContainer, false))
	cmd.Command("rebuild").Flags(launcher).RunFunc(launcher.InstanceFunc(launcher.Rebuild, true))
	cmd.Command("rename").Flags(launcher).RunFunc(launcher.Rename)
	cmd.Command("create-devices").Flags(launcher).RunFunc(launcher.InstanceFunc(launcher.CreateDevices, true))
	cmd.Command("create-profile").Flags(launcher).RunFunc(launcher.InstanceFunc(launcher.CreateProfile, false))

	var info cli.InfoOps
	cmd.Command("ostypes").RunFunc(info.ListOSTypes)

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

	configCmd := cmd.Command("config")
	parse := &ParseOp{}
	configCmd.Command("parse").Flags(parse).RunFunc(parse.Parse)
	configCmd.Command("print").Flags(parse).RunFunc(parse.Print)
	configOps := &ConfigOps{}
	configCmd.Command("formats").RunFunc(configOps.Formats)
	configCmd.Command("properties").RunFunc(configOps.PrintProperties)
	configCmd.Command("includes").RunFunc(configOps.Includes)
	configCmd.Command("script").RunFunc(configOps.Script)

	containerOps := &cli.InstanceOps{Client: client}
	containerCmd := cmd.Command("instance")
	containerCmd.Flags(containerOps)
	containerCmd.Command("profiles").RunFunc(containerOps.Profiles)
	containerCmd.Command("wait").RunFunc(containerOps.Wait)
	containerCmd.Command("devices").RunFunc(containerOps.Devices)
	containerCmd.Command("hwaddr").RunFunc(containerOps.ListHwaddr)
	containerCmd.Command("publish").RunFunc(containerOps.PublishInstance)

	cloudconfig := &cli.Cloudconfig{Client: client}
	cloudconfigInstanceOps := &cli.CloudconfigInstanceOps{Cloudconfig: cloudconfig}
	containerCmd.Command("cloudconfig").Flags(cloudconfigInstanceOps).RunFunc(cloudconfigInstanceOps.Apply)

	cloudconfigOps := &cli.CloudconfigOps{Cloudconfig: cloudconfig}
	cmd.Command("cloudconfig").Flags(cloudconfigOps).RunFunc(cloudconfigOps.Apply)

	imageCmd := cmd.Command("image")
	imageCmd.Command("instances").Flags(containerOps).RunFunc(containerOps.ListImages)

	networkOp := &cli.NetworkOp{Client: client}
	containerCmd.Command("addresses").Flags(networkOp).RunFunc(networkOp.ExportAddresses)

	numberOp := &cli.AssignNumbers{Client: client}
	containerCmd.Command("number").Flags(numberOp).RunFunc(numberOp.Run)

	exportOps := &ExportOps{Client: client}
	cmd.Command("export").Flags(exportOps).RunFunc(exportOps.Export)
	cmd.Command("import").Flags(exportOps).RunFunc(exportOps.Import)

	var migrate Migrate
	cmd.Command("copy-filesystems").Flags(&migrate).RunFunc(migrate.CopyFilesystems)

	usage.Apply(&cmd, usageData)

	return &cmd
}
