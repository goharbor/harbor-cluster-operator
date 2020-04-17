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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// HarborClusterSpec defines the desired state of HarborCluster
type HarborClusterSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of HarborCluster. Edit HarborCluster_types.go to remove/update
	Foo string `json:"foo,omitempty"`
}

// HarborClusterStatus defines the observed state of HarborCluster
type HarborClusterStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true

// HarborCluster is the Schema for the harborclusters API
type HarborCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   HarborClusterSpec   `json:"spec,omitempty"`
	Status HarborClusterStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// HarborClusterList contains a list of HarborCluster
type HarborClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []HarborCluster `json:"items"`
}

func init() {
	SchemeBuilder.Register(&HarborCluster{}, &HarborClusterList{})
}

type CRStatus struct {
	Phase           Phase  `json:"phase,omitempty"`
	ExternalService string `json:"service,omitempty"`
	AvailableNodes  int    `json:"availableNodes,omitempty"`
}

type Phase string

const (
	PendingPhase    Phase = "Pending"
	DeployingPhase  Phase = "Deploying"
	ReadyPhase      Phase = "Ready"
	UpgradingPhase  Phase = "Upgrading"
	DestroyingPhase Phase = "Destroying"
)
