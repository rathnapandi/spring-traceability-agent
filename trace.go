package main

import (
	"os"

	"git.ecd.axway.org/apigov/apic_agents_sdk/pkg/config"
	"git.ecd.axway.org/apigov/apic_agents_sdk/pkg/util/errors"

	_ "git.ecd.axway.org/apigov/apic_agents_sdk/pkg/traceability"
	agentCmd "github.com/rathnapandi/spring-traceability-agent/cmd"

	"github.com/elastic/beats/v7/filebeat/cmd"
)

func main() {
	err := agentCmd.InitEnvFileFlag()
	if err != nil {
		wrappedErr := errors.Wrap(config.ErrEnvConfigOverride, err.Error())
		cmd.RootCmd.Println("Error:", wrappedErr.Error())
		cmd.RootCmd.Println(cmd.RootCmd.UsageString())
		os.Exit(1)
	}

	if err := cmd.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
