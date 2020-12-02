/*
 * API Server specification.
 *
 * No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)
 *
 * API version: SNAPSHOT
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package v1alpha1

// ConsumerInstanceSpecIcon Image for the Catalog Item. If not present, the icon on the APISevice will be used in the Catalog Item.
type ConsumerInstanceSpecIcon struct {
	// Content-Type of the image.
	ContentType string `json:"contentType,omitempty"`
	// Base64 encoded image.
	Data string `json:"data,omitempty"`
}