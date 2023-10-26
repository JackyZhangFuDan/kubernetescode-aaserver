package v1alpha1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

type DbVolume string

const DbVolumeBig DbVolume = "BIG"
const DbVolumeSmall DbVolume = "SMALL"
const DbVolumeMedium DbVolume = "MEDIUM"

type ProvisionRequestSpec struct {
	IngressEntrance  string   `json:"ingressEntrance" `
	BusinessDbVolume DbVolume `json:"businessDbVolume"`
	NamespaceName    string   `json:"namespaceName"`
}

type ProvisionRequestStatus struct {
	IngressReady bool `json:"ingressReady" `
	DbReady      bool `json:"dbReady"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type ProvisionRequest struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" `

	Spec   ProvisionRequestSpec   `json:"spec,omitempty"`
	Status ProvisionRequestStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type ProvisionRequestList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty" `

	Items []ProvisionRequest `json:"items"`
}
