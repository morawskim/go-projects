/*
Copyright 2024.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// TraefikSpec defines the desired state of Traefik
type TraefikSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	LookForLabel         string `json:"lookForLabel,omitempty"`
	TargetNamespace      string `json:"targetNamespace,omitempty"`
	TargetDeploymentName string `json:"targetDeploymentName,omitempty"`
	TargetConfigMapName  string `json:"targetConfigMapName,omitempty"`
}

// TraefikStatus defines the observed state of Traefik
type TraefikStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// Traefik is the Schema for the traefiks API
type Traefik struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TraefikSpec   `json:"spec,omitempty"`
	Status TraefikStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// TraefikList contains a list of Traefik
type TraefikList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Traefik `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Traefik{}, &TraefikList{})
}
