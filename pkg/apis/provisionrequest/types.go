package provisionrequest

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

type ProvisionRequest struct {
	metav1.TypeMeta
	metav1.ObjectMeta

	Spec   ProvisionRequestSpec
	Status ProvisionRequestStatus
}

type ProvisionRequestList struct {
	metav1.TypeMeta
	metav1.ListMeta

	Items []ProvisionRequest
}
