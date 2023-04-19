/*
Copyright 2023 wu123.

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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// TraefikServiceSpec defines the desired state of TraefikService
type TraefikServiceSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of TraefikService. Edit traefikservice_types.go to remove/update
	Entrypoint string `json:"entrypoint"`
	URL        string `json:"url"`
}

// TraefikServiceStatus defines the observed state of TraefikService
type TraefikServiceStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Active         bool         `json:"active"`
	LastUpdateTime *metav1.Time `json:"lastupdatetime"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// TraefikService is the Schema for the traefikservices API
type TraefikService struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TraefikServiceSpec   `json:"spec,omitempty"`
	Status TraefikServiceStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// TraefikServiceList contains a list of TraefikService
type TraefikServiceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []TraefikService `json:"items"`
}

func init() {
	SchemeBuilder.Register(&TraefikService{}, &TraefikServiceList{})
}
