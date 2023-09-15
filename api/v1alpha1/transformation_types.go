/*
Copyright 2023 WatcherWhale.

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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type TransformationSpec struct {
	Target   TargetSpec        `json:"target"`
	Template map[string]string `json:"template"`
	Sources  []SourceSpec      `json:"sources"`
}

type TargetSpec struct {
	Name        string            `json:"name"`
	Kind        string            `json:"kind"`
	Labels      map[string]string `json:"labels"`
	Annotations map[string]string `json:"annotations"`
}

type SourceSpec struct {
	Kind string `json:"kind"`
	Name string `json:"name"`
}

// TransformationStatus defines the observed state of Transformation
type TransformationStatus struct {
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Transformation is the Schema for the transformations API
type Transformation struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TransformationSpec   `json:"spec,omitempty"`
	Status TransformationStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// TransformationList contains a list of Transformation
type TransformationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Transformation `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Transformation{}, &TransformationList{})
}
