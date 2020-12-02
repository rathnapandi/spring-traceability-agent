package config

import "git.ecd.axway.org/apigov/apic_agents_sdk/pkg/util/errors"

// Config errors
var (
	ErrGatewayConfig = errors.Newf(3500, "error apigateway.getHeaders is set to true and apigateway.%s not set in config")
)
