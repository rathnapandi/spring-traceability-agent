package main

import (
	"flag"
	"fmt"
	"git.ecd.axway.org/apigov/apic_agents_sdk/pkg/config"
	"git.ecd.axway.org/apigov/apic_agents_sdk/pkg/util"
	"git.ecd.axway.org/apigov/apic_agents_sdk/pkg/util/errors"
	"os"

	corecmd "git.ecd.axway.org/apigov/apic_agents_sdk/pkg/cmd"
	"git.ecd.axway.org/apigov/apic_agents_sdk/pkg/cmd/service"
	hc "git.ecd.axway.org/apigov/apic_agents_sdk/pkg/util/healthcheck"

	//statuscmd "git.ecd.axway.org/apigov/v7_traceability_agent/pkg/cmd/status"

	"github.com/elastic/beats/v7/filebeat/cmd"
	libcmd "github.com/elastic/beats/v7/libbeat/cmd"
	//libcmd "github.com/elastic/beats/v7/libbeat/cmd"
	"github.com/elastic/beats/v7/libbeat/cmd/instance"
	"github.com/rathnapandi/spring-traceability-agent/pkg/beater"
	"github.com/spf13/pflag"
)

// Name of this beat
var Name = "traceability_agent"

const modules = "modules"

func main() {


	Init()
	err := InitEnvFileFlag()
	if err != nil {
		wrappedErr := errors.Wrap(config.ErrEnvConfigOverride, err.Error())
		cmd.RootCmd.Println("Error:", wrappedErr.Error())
		cmd.RootCmd.Println(cmd.RootCmd.UsageString())
		os.Exit(1)
	}

	//if err := cmd.RootCmd.Execute(); err != nil {
	//	//agent.UpdateStatus(agent.AgentFailed, err.Error())
	//	os.Exit(1)
	//}
	//agent.UpdateStatus(agent.AgentStopped, "")
}


func Init() {
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
	hc.SetNameAndVersion(Name, version)
	println(cmd.RootCmd)
}


// InitEnvFileFlag - Initialize the envFile Flag
func InitEnvFileFlag() error {
	flag.String(corecmd.EnvFileFlag, "", corecmd.EnvFileFlagDesciption)
	cmd.RootCmd.PersistentFlags().AddGoFlag(flag.CommandLine.Lookup(corecmd.EnvFileFlag))
	flag.Parse()
	fmt.Println("init")
	envFile, _ := cmd.RootCmd.PersistentFlags().GetString(corecmd.EnvFileFlag)
	err := util.LoadEnvFromFile(envFile)
	if err != nil {
		return err
	}
	return nil
}
