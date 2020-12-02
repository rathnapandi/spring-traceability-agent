/*
 * API Server specification.
 *
 * No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)
 *
 * API version: SNAPSHOT
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package v1alpha1

// EdgeDataplaneSpecApiGateway Axway API Gateway configuration.
type EdgeDataplaneSpecApiGateway struct {
	// Host name where Axway API Gateway is deployed
	Host string `json:"host,omitempty"`
	// Interval the Agent will poll API Gateway. Defaults to '1m' indicating 1 minute
	PollInterval string `json:"pollInterval,omitempty"`
	// API Gateway Admin port. Defaults to 8090
	Port int32 `json:"port,omitempty"`
}