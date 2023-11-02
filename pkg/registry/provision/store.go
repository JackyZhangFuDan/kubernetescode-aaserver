package provision

import (
	"github.com/kubernetescode-aaserver/pkg/apis/provision"
	"github.com/kubernetescode-aaserver/pkg/registry"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/registry/generic"
	gRegistry "k8s.io/apiserver/pkg/registry/generic/registry"
	"k8s.io/apiserver/pkg/registry/rest"
)

func NewREST(scheme *runtime.Scheme, optsGetter generic.RESTOptionsGetter) (*registry.REST, *registry.StatusREST, error) {
	strategy := NewStrategy(scheme)

	store := &gRegistry.Store{
		NewFunc:                   func() runtime.Object { return &provision.ProvisionRequest{} },
		NewListFunc:               func() runtime.Object { return &provision.ProvisionRequestList{} },
		PredicateFunc:             MatchService,
		DefaultQualifiedResource:  provision.Resource("provisionrequests"),
		SingularQualifiedResource: provision.Resource("provisionrequest"),

		CreateStrategy:      strategy,
		UpdateStrategy:      strategy,
		DeleteStrategy:      strategy,
		ResetFieldsStrategy: strategy,

		TableConvertor: rest.NewDefaultTableConvertor(provision.Resource("provisionrequests")),
	}
	options := &generic.StoreOptions{RESTOptions: optsGetter, AttrFunc: GetAttrs}
	if err := store.CompleteWithOptions(options); err != nil {
		return nil, nil, err
	}

	statusStrategy := NewStatusStrategy(strategy)
	statusStore := *store
	statusStore.UpdateStrategy = statusStrategy
	statusStore.ResetFieldsStrategy = statusStrategy

	return &registry.REST{Store: store}, &registry.StatusREST{Store: &statusStore}, nil
}
