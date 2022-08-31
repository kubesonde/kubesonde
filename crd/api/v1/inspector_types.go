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

type ActionType string

const (
	ALLOW ActionType = "Allow"
	DENY  ActionType = "Deny"
)

type ProbingAction struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Action is the expected outcome of the probe. It might have values "allow" or "deny"
	Action ActionType `json:"action,omitempty"`
	// FromPodSelector is a selector for the origin Pod or a set of pods
	FromPodSelector string `json:"fromPodSelector,omitempty"`
	// ToPodSelector is a selector for the destination Pod or a set of pods
	// +optional
	ToPodSelector string `json:"toPodSelector,omitempty"`
	// Port is the probing port for ToPodSelector defaults to 80
	// +optional
	Port string `json:"port,omitempty"`
	// Protocol is the protocol to use when probing ToPodSelector
	// +optional
	Protocol string `json:"protocol,omitempty"`
	// Endpoint to probe
	// +optional
	Endpoint string `json:"endpoint,omitempty"`
	// Url to probe
	// +optional
	Url string `json:"url,omitempty"`
}

// KubesondeSpec defines the desired state of Kubesonde
type KubesondeSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Namespace indicates the target namespace for the probe
	// +optional
	Namespace string `json:"namespace,omitempty"`
	// Actions is the set of rules to be tested
	// +optional
	Actions []ProbingAction `json:"actions,omitempty"`
}

// KubesondeStatus defines the observed state of Kubesonde
type KubesondeStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	// TODO: Add information regarding the running containers

	// Information when was the last time the probe was run.
	// +optional
	LastProbeTime *metav1.Time `json:"lastProbeTime,omitempty"`
}

// +kubebuilder:object:root=true

// Kubesonde is the Schema for the Kubesondes API
type Kubesonde struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   KubesondeSpec   `json:"spec,omitempty"`
	Status KubesondeStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// KubesondeList contains a list of Kubesonde
type KubesondeList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Kubesonde `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Kubesonde{}, &KubesondeList{})
}
