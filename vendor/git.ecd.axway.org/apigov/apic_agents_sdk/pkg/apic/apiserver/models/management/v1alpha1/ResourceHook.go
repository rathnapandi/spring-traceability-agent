/*
 * This file is automatically generated
 */

package v1alpha1

import (
	"encoding/json"

	apiv1 "git.ecd.axway.org/apigov/apic_agents_sdk/pkg/apic/apiserver/models/api/v1"
)

var (
	_ResourceHookGVK = apiv1.GroupVersionKind{
		GroupKind: apiv1.GroupKind{
			Group: "management",
			Kind:  "ResourceHook",
		},
		APIVersion: "v1alpha1",
	}
)

const (
	ResourceHookScope = "Integration"

	ResourceHookResource = "resourcehooks"
)

func ResourceHookGVK() apiv1.GroupVersionKind {
	return _ResourceHookGVK
}

func init() {
	apiv1.RegisterGVK(_ResourceHookGVK, ResourceHookScope, ResourceHookResource)
}

// ResourceHook Resource
type ResourceHook struct {
	apiv1.ResourceMeta

	Spec ResourceHookSpec `json:"spec"`
}

// FromInstance converts a ResourceInstance to a ResourceHook
func (res *ResourceHook) FromInstance(ri *apiv1.ResourceInstance) error {
	m, err := json.Marshal(ri.Spec)
	if err != nil {
		return err
	}

	spec := &ResourceHookSpec{}
	err = json.Unmarshal(m, spec)
	if err != nil {
		return err
	}

	*res = ResourceHook{ResourceMeta: ri.ResourceMeta, Spec: *spec}

	return err
}

// AsInstance converts a ResourceHook to a ResourceInstance
func (res *ResourceHook) AsInstance() (*apiv1.ResourceInstance, error) {
	m, err := json.Marshal(res.Spec)
	if err != nil {
		return nil, err
	}

	spec := map[string]interface{}{}
	err = json.Unmarshal(m, &spec)
	if err != nil {
		return nil, err
	}

	return &apiv1.ResourceInstance{ResourceMeta: res.ResourceMeta, Spec: spec}, nil
}
