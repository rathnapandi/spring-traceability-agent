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
// CatalogItemInit struct for CatalogItemInit
type CatalogItemInit struct {
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
	// Description of the catalog item
	Description string `json:"description,omitempty"`
	Properties []CatalogItemProperty `json:"properties"`
	Tags []string `json:"tags,omitempty"`
	Visibility string `json:"visibility"`
	Subscription CatalogItemSubscriptionDefinition `json:"subscription"`
	Revision CatalogItemInitRevision `json:"revision"`
	// A list of categories ids.
	CategoriesReferences []string `json:"categoriesReferences,omitempty"`
}