/*
Copyright 2024.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either expzress or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// SnowflakeResourceSpec defines the desired state of SnowflakeResource
type SnowflakeResourceSpec struct {
	// ObjectType is the type of the Snowflake object
	// +kubebuilder:validation:Required
	ObjectType string `json:"objectType"`

	// ObjectConfig is a list of parameters for the Snowflake object
	// +kubebuilder:validation:Required
	ObjectConfig []Param `json:"objectConfig"`
}

// Param defines the parameters for the snowflake object
type Param struct {
	// Key is the name of the parameter
	// +kubebuilder:validation:Required
	Key string `json:"key"`

	// Value is the value of the parameter
	// +kubebuilder:validation:Required
	Value string `json:"value"`
}

// SnowflakeResourceStatus defines the observed state of SnowflakeResource
type SnowflakeResourceStatus struct {
	Ready      bool   `json:"ready"`
	ObjectType string `json:"type"`
	ObjectName string `json:"name"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// SnowflakeResource is the Schema for the snowflakeresources API
type SnowflakeResource struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SnowflakeResourceSpec   `json:"spec,omitempty"`
	Status SnowflakeResourceStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// SnowflakeResourceList contains a list of SnowflakeResource
type SnowflakeResourceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SnowflakeResource `json:"items"`
}

func init() {
	SchemeBuilder.Register(&SnowflakeResource{}, &SnowflakeResourceList{})
}
