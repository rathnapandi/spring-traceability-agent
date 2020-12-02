/*
 * API Server specification.
 *
 * No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)
 *
 * API version: SNAPSHOT
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package v1alpha1

// AwsDiscoveryAgentState struct for AwsDiscoveryAgentState
type AwsDiscoveryAgentState struct {
	// A way to communicate details about the current status by the agent
	Description string `json:"description,omitempty"`
	// Agent status:  * waiting - Default status to indicate that resource is defined, but agent have not connected yet  * running - Passed all health checks.  Up and running  * stopped - Failed health checks.
	Status string `json:"status,omitempty"`
	// Version name for the agent revision.
	Version string `json:"version,omitempty"`
}
