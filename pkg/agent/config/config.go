package config

import (
	"time"

	v1 "git.ecd.axway.org/apigov/apic_agents_sdk/pkg/apic/apiserver/models/api/v1"
	"git.ecd.axway.org/apigov/apic_agents_sdk/pkg/apic/apiserver/models/management/v1alpha1"
	corecfg "git.ecd.axway.org/apigov/apic_agents_sdk/pkg/config"
	"git.ecd.axway.org/apigov/apic_agents_sdk/pkg/util/exception"

)

// AgentConfig - Agent Configuration
type AgentConfig struct {
	Central corecfg.CentralConfig             `config:"central"`
	//Gateway APIGatewayConfiguration           `config:"apigateway"`
	//Manager apimgrcfg.APIManagerConfiguration `config:"apimanager"`
	Status  corecfg.StatusConfig              `config:"status"`
}

// New - Creates a new Config object with defaults
func New() AgentConfig {
	return AgentConfig{
		Central: corecfg.NewCentralConfig(corecfg.TraceabilityAgent),
		//Gateway: NewGatewayConfig(),
		//Manager: apimgrcfg.APIManagerConfiguration{
		//	PollInterval: 1 * time.Minute,
		//	TLS:          corecfg.NewTLSConfig(),
		//},
		Status: corecfg.NewStatusConfig(),
	}
}

// APIGatewayConfiguration - APIGateway Configuration
type APIGatewayConfiguration struct {
	corecfg.IConfigValidator
	corecfg.IResourceConfigCallback
	Host           string            `config:"host"`
	Port           int               `config:"port"`
	User           string            `config:"auth.username"`
	Password       string            `config:"auth.password"`
	EnableAPICalls bool              `config:"getHeaders"`
	PollInterval   time.Duration     `config:"pollInterval"`
	TLS            corecfg.TLSConfig `config:"ssl"`
	ProxyURL       string            `config:"proxyUrl"`
}

// NewGatewayConfig -
func NewGatewayConfig() APIGatewayConfiguration {
	return APIGatewayConfiguration{
		EnableAPICalls: true,
		PollInterval:   time.Minute,
		TLS:            corecfg.NewTLSConfig(),
	}
}

// ValidateCfg - Validates the config, implementing IConfigInterface
func (c *APIGatewayConfiguration) ValidateCfg() (err error) {
	exception.Block{
		Try: func() {
			c.validateConfig()
		},
		Catch: func(e error) {
			err = e
		},
	}.Do()

	return
}

func (c *APIGatewayConfiguration) validateConfig() {
	// If API calls are disabled these config options are not required
	if c.EnableAPICalls {
		if c.Host == "" {
			exception.Throw(ErrGatewayConfig.FormatError("host"))
		}

		if c.Port == 0 {
			exception.Throw(ErrGatewayConfig.FormatError("port"))
		}

		if c.User == "" {
			exception.Throw(ErrGatewayConfig.FormatError("auth.username"))
		}

		if c.Password == "" {
			exception.Throw(ErrGatewayConfig.FormatError("auth.password"))
		}
	}

	if c.PollInterval == 0 {
		c.PollInterval = time.Minute
	}
}

// ParseAPIManagerConfig - parse the props and create an API Manager Configuration structure
func createFromResources(dp *v1alpha1.EdgeDataplane, da *v1alpha1.EdgeTraceabilityAgent) (*APIGatewayConfiguration, error) {

	cfg := &APIGatewayConfiguration{
		Host:           dp.Spec.ApiGatewayManager.Host,
		Port:           int(dp.Spec.ApiGatewayManager.Port),
		EnableAPICalls: da.Spec.Config.ProcessHeaders,
		PollInterval:   1 * time.Minute,
	}

	if dp.Spec.ApiGatewayManager.PollInterval != "" {
		resCfgPollInterval, err := time.ParseDuration(dp.Spec.ApiGatewayManager.PollInterval)
		if err != nil {
			return nil, err
		}
		cfg.PollInterval = resCfgPollInterval
	}
	return cfg, nil
}

// ApplyResources - Applies the agent and dataplane resource to config
func (c *APIGatewayConfiguration) ApplyResources(dataplaneResource *v1.ResourceInstance, agentResource *v1.ResourceInstance) error {
	dp := &v1alpha1.EdgeDataplane{}
	err := dp.FromInstance(dataplaneResource)
	if err != nil {
		return err
	}

	da := &v1alpha1.EdgeTraceabilityAgent{}
	err = da.FromInstance(agentResource)
	if err != nil {
		return err
	}

	cfgFromRes, err := createFromResources(dp, da)
	if err != nil {
		return err
	}
	// Check if local config not set to default and resource is non default use the config from res
	if c.Host == "localhost" && cfgFromRes.Host != "localhost" {
		c.Host = cfgFromRes.Host
	}
	if c.Port == 8090 && cfgFromRes.Port != 8090 {
		c.Port = cfgFromRes.Port
	}
	if c.PollInterval == 1*time.Minute && cfgFromRes.PollInterval != 1*time.Minute {
		c.PollInterval = cfgFromRes.PollInterval
	}
	if c.EnableAPICalls && !cfgFromRes.EnableAPICalls {
		c.EnableAPICalls = cfgFromRes.EnableAPICalls
	}
	return nil
}
