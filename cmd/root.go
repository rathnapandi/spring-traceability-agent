package cmd

import (
	"flag"
	"fmt"

	"git.ecd.axway.org/apigov/apic_agents_sdk/pkg/util"

	corecmd "git.ecd.axway.org/apigov/apic_agents_sdk/pkg/cmd"
	"git.ecd.axway.org/apigov/apic_agents_sdk/pkg/cmd/service"

	//statuscmd "git.ecd.axway.org/apigov/v7_traceability_agent/pkg/cmd/status"

	"github.com/elastic/beats/v7/filebeat/cmd"
	libcmd "github.com/elastic/beats/v7/libbeat/cmd"
	"github.com/elastic/beats/v7/libbeat/cmd/instance"
	"github.com/spf13/pflag"

	"github.com/rathnapandi/spring-traceability-agent/pkg/beater"
)

// Name of this beat
var Name = "traceability_agent"

const modules = "modules"

func init() {
	// RootCmd to handle beats cli
	var runFlags = pflag.NewFlagSet(Name, pflag.ExitOnError)
	runFlags.AddGoFlag(flag.CommandLine.Lookup("once"))
	runFlags.AddGoFlag(flag.CommandLine.Lookup(modules))
	cmd.Name = Name
	version := fmt.Sprintf("%s-%s", corecmd.BuildVersion, corecmd.BuildCommitSha)
	cmdSettings := instance.Settings{RunFlags: runFlags, Name: cmd.Name, Version: version}
	cmd.RootCmd = libcmd.GenRootCmdWithSettings(beater.New, cmdSettings)

	// Add the Status subcommand
	//cmd.RootCmd.AddCommand(statuscmd.GenStatusCmd(cmdSettings, beater.New))

	// Add the Service subcommand
	cmd.RootCmd.AddCommand(service.GenServiceCmd("path.config"))

	cmd.RootCmd.PersistentFlags().AddGoFlag(flag.CommandLine.Lookup("M"))
	cmd.RootCmd.TestCmd.Flags().AddGoFlag(flag.CommandLine.Lookup(modules))
	cmd.RootCmd.SetupCmd.Flags().AddGoFlag(flag.CommandLine.Lookup(modules))

	// set the healthcheck values
	// hc.SetNameAndVersion(Name, version)
}

// InitEnvFileFlag - Initialize the envFile Flag
func InitEnvFileFlag() error {
	flag.String(corecmd.EnvFileFlag, "", corecmd.EnvFileFlagDesciption)
	cmd.RootCmd.PersistentFlags().AddGoFlag(flag.CommandLine.Lookup(corecmd.EnvFileFlag))
	flag.Parse()
	envFile, _ := cmd.RootCmd.PersistentFlags().GetString(corecmd.EnvFileFlag)
	err := util.LoadEnvFromFile(envFile)
	if err != nil {
		return err
	}
	return nil
}
