/*


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

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

type Incident struct {
	Name         string        `json:"name,omitempty"`
	Remediations []Remediation `json:"remediations"`
}

type Remediation struct {
	Name string `json:"active,omitempty"`
}

// +kubebuilder:printcolumn:JSONPath=".spec.Active",name=Active,type=bool
type PloverSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	Active    bool       `json:"active,omitempty"`
	Incidents []Incident `json:"incidents"`
}

// PloverStatus defines the observed state of Plover
type PloverStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true

// Plover is the Schema for the plovers API
type Plover struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PloverSpec   `json:"spec,omitempty"`
	Status PloverStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// PloverList contains a list of Plover
type PloverList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Plover `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Plover{}, &PloverList{})
}
