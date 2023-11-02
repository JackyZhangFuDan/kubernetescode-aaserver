package registry

import (
	"context"
	"fmt"

	"github.com/kubernetescode-aaserver/pkg/apis/provision"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	gRegistry "k8s.io/apiserver/pkg/registry/generic/registry"
	"k8s.io/apiserver/pkg/registry/rest"
	"sigs.k8s.io/structured-merge-diff/v4/fieldpath"
)

type REST struct {
	*gRegistry.Store
}

func RESTWithErrorHandler(storage rest.StandardStorage, err error) rest.StandardStorage {
	if err != nil {
		err = fmt.Errorf("fail when creating rest storage for a resource due to %v, will die", err)
		panic(err)
	}
	return storage
}

// implements the ShortNamesProvider interface.
func (r *REST) ShortNames() []string {
	return []string{"pr"}
}

type StatusREST struct {
	*gRegistry.Store
}

func (r *StatusREST) New() runtime.Object {
	return &provision.ProvisionRequest{}
}

// Destroy cleans up resources on shutdown.
func (r *StatusREST) Destroy() {
	// Given that underlying store is shared with REST,
	// we don't destroy it here explicitly.
}

// Get retrieves the object from the storage. It is required to support Patch.
func (r *StatusREST) Get(ctx context.Context, name string, options *metav1.GetOptions) (runtime.Object, error) {
	return r.Store.Get(ctx, name, options)
}

func (r *StatusREST) Update(ctx context.Context, name string, objInfo rest.UpdatedObjectInfo, createValidation rest.ValidateObjectFunc, updateValidation rest.ValidateObjectUpdateFunc, forceAllowCreate bool, options *metav1.UpdateOptions) (runtime.Object, bool, error) {
	// We are explicitly setting forceAllowCreate to false in the call to the underlying storage because
	// subresources should never allow create on update.
	return r.Store.Update(ctx, name, objInfo, createValidation, updateValidation, false, options)
}

// implements rest.ResetFieldsStrategy
func (r *StatusREST) GetResetFields() map[fieldpath.APIVersion]*fieldpath.Set {
	return r.Store.GetResetFields()
}

func (r *StatusREST) ConvertToTable(ctx context.Context, object runtime.Object, tableOptions runtime.Object) (*metav1.Table, error) {
	return r.Store.ConvertToTable(ctx, object, tableOptions)
}
