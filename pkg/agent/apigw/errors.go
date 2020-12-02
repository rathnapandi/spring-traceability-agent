package apigw

import "git.ecd.axway.org/apigov/apic_agents_sdk/pkg/util/errors"

// API Gateway errors
var (
	// Healthcheck
	ErrV7HealthcheckHost = errors.Newf(3501, "%s Failed. Error communicating with API Gateway: %s. Check API Gateway host configuration values")
	ErrV7HealthcheckAuth = errors.Newf(3502, "%s Failed. Error sending request to API Gateway. HTTP response code %v. Check API Gateway authentication configuration values")

	// Event Processing
	ErrEventNoMsg        = errors.Newf(3510, "the log event had no message field: %s")
	ErrEventMsgStructure = errors.Newf(3511, "could not parse the log event: %s")
	ErrTrxnDataGet       = errors.New(3512, "could not retrieve the transaction data from API Gateway")
	ErrTrxnDataProcess   = errors.New(3513, "could not process the transaction data")
	ErrTrxnHeaders       = errors.Newf(3514, "could not process the transaction headers: %s")
	ErrProtocolStructure = errors.Newf(3515, "could not parse the %s transaction details: %s")
	ErrCreateCondorEvent = errors.New(3516, "error creating the AMPLIFY Visibility event")

	// API Gateway Communication
	ErrAPIGWRequest  = errors.New(3530, "error encountered sending a request to API Gateway")
	ErrAPIGWResponse = errors.Newf(3531, "unexpected HTTP response code %s in response from API Gateway")
)
