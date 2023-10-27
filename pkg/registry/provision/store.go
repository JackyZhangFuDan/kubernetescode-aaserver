package provision

import (
	"github.com/kubernetescode-aaserver/pkg/apis/provision"
	"github.com/kubernetescode-aaserver/pkg/registry"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/registry/generic"
	gRegistry "k8s.io/apiserver/pkg/registry/generic/registry"
	"k8s.io/apiserver/pkg/registry/rest"
)

func NewREST(scheme *runtime.Scheme, optsGetter generic.RESTOptionsGetter) (*registry.REST, error) {
	strategy := NewStrategy(scheme)

	store := &gRegistry.Store{
		NewFunc:                  func() runtime.Object { return &provision.ProvisionRequest{} },
		NewListFunc:              func() runtime.Object { return &provision.ProvisionRequestList{} },
		PredicateFunc:            MatchJenkinsService,
		DefaultQualifiedResource: provision.Resource("provisionrequests"),

		CreateStrategy: strategy,
		UpdateStrategy: nil,
		DeleteStrategy: nil,

		TableConvertor: rest.NewDefaultTableConvertor(provision.Resource("provisionrequests")),
	}
	options := &generic.StoreOptions{RESTOptions: optsGetter, AttrFunc: GetAttrs}
	if err := store.CompleteWithOptions(options); err != nil {
		return nil, err
	}
	return &registry.REST{Store: store}, nil
}
