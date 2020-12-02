package beater

import (
	"fmt"
	//"github.com/elastic/beats/x-pack/elastic-agent/pkg/config"

	"github.com/elastic/beats/v7/filebeat/beater"
	"github.com/elastic/beats/v7/libbeat/beat"
	"github.com/elastic/beats/v7/libbeat/common"

	"git.ecd.axway.org/apigov/apic_agents_sdk/pkg/agent"
	"git.ecd.axway.org/apigov/apic_agents_sdk/pkg/apic"
	"github.com/rathnapandi/spring-traceability-agent/pkg/agent/apigw"
	"github.com/rathnapandi/spring-traceability-agent/pkg/agent/config"

	corecfg "git.ecd.axway.org/apigov/apic_agents_sdk/pkg/config"
	"git.ecd.axway.org/apigov/apic_agents_sdk/pkg/traceability"
	coreerrors "git.ecd.axway.org/apigov/apic_agents_sdk/pkg/util/errors"
	hc "git.ecd.axway.org/apigov/apic_agents_sdk/pkg/util/healthcheck"
)

// AgentCfg -
var AgentCfg config.AgentConfig

// GetConfig -
func GetConfig() *config.AgentConfig {
	return &AgentCfg
}

// New creates an instance of aws_apigw_traceability_agent.
func New(b *beat.Beat, cfg *common.Config) (beat.Beater, error) {
	fmt.Println("beater new")
	AgentCfg := config.New()
	if err := cfg.Unpack(&AgentCfg); err != nil {
		return nil, fmt.Errorf("Error reading config file: %v", err)
	}
	err := agent.Initialize(AgentCfg.Central)
	if err != nil {
		return nil, err
	}
	// Merge agent config with agent agent resource

	// Validate agent config
	err = corecfg.ValidateConfig(AgentCfg)
	if err != nil {
		return nil, err
	}

	// Create APIC Client to register the healthcheck
	apic.New(AgentCfg.Central)

	// Init the healthcheck API
	hc.SetStatusConfig(AgentCfg.Status)
	hc.HandleRequests()

	//// Start watching for API Manager changes
	//_, err = apimanager.New(AgentCfg.Manager)
	//if err != nil {
	//	agent.UpdateStatus(agent.AgentFailed, err.Error())
	//	return nil, err
	//}

	// Now that we have the config, set the proxy if necessary.
	err = AgentCfg.Central.SetProxyEnvironmentVariable()
	if err != nil {
		return nil, coreerrors.Wrap(traceability.ErrSettingProxy, err.Error())
	}

	if hc.RunChecks() != hc.OK {
		return nil, coreerrors.ErrInitServicesNotReady
	}

	eventProcessor := apigw.New()
	traceability.SetOutputEventProcessor(eventProcessor)

	// Initialize the filebeat to read events
	return beater.New(b, cfg)
}
