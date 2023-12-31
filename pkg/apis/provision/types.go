package provision

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

type DbVolume string

const DbVolumeBig DbVolume = "BIG"
const DbVolumeSmall DbVolume = "SMALL"
const DbVolumeMedium DbVolume = "MEDIUM"

type ProvisionRequestSpec struct {
	IngressEntrance  string
	BusinessDbVolume DbVolume
	NamespaceName    string
}

type ProvisionRequestStatus struct {
	IngressReady bool
	DbReady      bool
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type ProvisionRequest struct {
	metav1.TypeMeta
	metav1.ObjectMeta

	Spec   ProvisionRequestSpec
	Status ProvisionRequestStatus
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type ProvisionRequestList struct {
	metav1.TypeMeta
	metav1.ListMeta

	Items []ProvisionRequest
}
