package install

import (
	provisionrequest "github.com/kubernetescode-aaserver/pkg/apis/provision"
	"github.com/kubernetescode-aaserver/pkg/apis/provision/v1alpha1"

	"k8s.io/apimachinery/pkg/runtime"
	util "k8s.io/apimachinery/pkg/util/runtime"
)

func Install(scheme *runtime.Scheme) {
	util.Must(provisionrequest.AddToScheme(scheme))
	util.Must(v1alpha1.AddToScheme(scheme))
	util.Must(scheme.SetVersionPriority(v1alpha1.SchemeGroupVersion))
}
