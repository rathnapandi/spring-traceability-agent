/*
 * Amplify Unified Catalog APIs
 *
 * APIs for Amplify Unified Catalog
 *
 * API version: 1.43.0
 * Contact: support@axway.com
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package unifiedcatalog
// CatalogItem struct for CatalogItem
type CatalogItem struct {
	// Generated identifier for the resource
	Id string `json:"id,omitempty"`
	// Id of the owning team of this catalog item
	OwningTeamId string `json:"owningTeamId,omitempty"`
	// Type of the definition for the catalog item
	DefinitionType string `json:"definitionType"`
	// Sub-Type of the definition for the catalog item
	DefinitionSubType string `json:"definitionSubType"`
	// Revision of the definition for the catalog item
	DefinitionRevision int32 `json:"definitionRevision"`
	// Name of the catalog item
	Name string `json:"name"`
	Relationships EntityRelationship `json:"relationships,omitempty"`
	// Description of the catalog item
	Description string `json:"description,omitempty"`
	Tags []string `json:"tags,omitempty"`
	Metadata AuditMetadata `json:"metadata,omitempty"`
	Visibility string `json:"visibility"`
	State string `json:"state"`
	Access string `json:"access,omitempty"`
	AvailableRevisions []int32 `json:"availableRevisions,omitempty"`
	// Latest version of the published revision.
	LatestVersion string `json:"latestVersion,omitempty"`
	// Number of subscriptions for the catalog item
	TotalSubscriptions int32 `json:"totalSubscriptions,omitempty"`
	LatestVersionDetails CatalogItemRevision `json:"latestVersionDetails,omitempty"`
	Image CatalogItemImage `json:"image,omitempty"`
	Categories []Category `json:"categories,omitempty"`
	Acl []AccessControlItem `json:"acl,omitempty"`
}
