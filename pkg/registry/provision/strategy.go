package provision

import (
	"context"
	"fmt"

	"github.com/kubernetescode-aaserver/pkg/apis/provision"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/storage"
	"k8s.io/apiserver/pkg/storage/names"
	"sigs.k8s.io/structured-merge-diff/v4/fieldpath"
)

type provisionRequestStrategy struct {
	runtime.ObjectTyper
	names.NameGenerator
}

func NewStrategy(typer runtime.ObjectTyper) provisionRequestStrategy {
	return provisionRequestStrategy{typer, names.SimpleNameGenerator}
}

func (provisionRequestStrategy) NamespaceScoped() bool {
	return true
}
func (provisionRequestStrategy) PrepareForCreate(ctx context.Context, obj runtime.Object) {

}
func (provisionRequestStrategy) Validate(ctx context.Context, obj runtime.Object) field.ErrorList {
	errs := field.ErrorList{}

	js := obj.(*provision.ProvisionRequest)
	if len(js.Spec.NamespaceName) == 0 {
		errs = append(errs,
			field.Required(field.NewPath("spec").Key("namespaceName"),
				"namespace name is required"))
	}
	if len(errs) > 0 {
		return errs
	} else {
		return nil
	}
}
func (provisionRequestStrategy) WarningsOnCreate(ctx context.Context, obj runtime.Object) []string {
	return []string{}
}
func (provisionRequestStrategy) Canonicalize(obj runtime.Object) {

}

func (provisionRequestStrategy) AllowUnconditionalUpdate() bool {
	return false
}
func (provisionRequestStrategy) PrepareForUpdate(ctx context.Context, obj runtime.Object, old runtime.Object) {

}
func (provisionRequestStrategy) ValidateUpdate(ctx context.Context, obj, old runtime.Object) field.ErrorList {
	return field.ErrorList{}
}
func (provisionRequestStrategy) WarningsOnUpdate(ctx context.Context, obj runtime.Object, old runtime.Object) []string {
	return []string{}
}
func (provisionRequestStrategy) AllowCreateOnUpdate() bool {
	return false
}
func (provisionRequestStrategy) GetResetFields() map[fieldpath.APIVersion]*fieldpath.Set {
	fields := map[fieldpath.APIVersion]*fieldpath.Set{
		"provision.mydomain.com/v1alpha1": fieldpath.NewSet(
			fieldpath.MakePathOrDie("status"),
		),
	}
	return fields
}

func GetAttrs(obj runtime.Object) (labels.Set, fields.Set, error) {
	object, ok := obj.(*provision.ProvisionRequest)
	if !ok {
		return nil, nil, fmt.Errorf("the object isn't a ProvisionRequest")
	}
	fs := generic.ObjectMetaFieldsSet(&object.ObjectMeta, true)
	return labels.Set(object.ObjectMeta.Labels), fs, nil
}

func MatchService(label labels.Selector, field fields.Selector) storage.SelectionPredicate {
	return storage.SelectionPredicate{
		Label:    label,
		Field:    field,
		GetAttrs: GetAttrs,
	}
}

type provisionRequestStatusStrategy struct {
	provisionRequestStrategy
}

func NewStatusStrategy(st provisionRequestStrategy) provisionRequestStatusStrategy {
	return provisionRequestStatusStrategy{st}
}

// GetResetFields returns the set of fields that get reset by the strategy
// and should not be modified by the user.
func (provisionRequestStatusStrategy) GetResetFields() map[fieldpath.APIVersion]*fieldpath.Set {
	return map[fieldpath.APIVersion]*fieldpath.Set{
		"provision.mydomain.com/v1alpha1": fieldpath.NewSet(
			fieldpath.MakePathOrDie("spec"),
			fieldpath.MakePathOrDie("metadata", "labels"),
		),
	}
}

// PrepareForUpdate clears fields that are not allowed to be set by end users on update of status
func (provisionRequestStatusStrategy) PrepareForUpdate(ctx context.Context, obj, old runtime.Object) {
	newObj := obj.(*provision.ProvisionRequest)
	oldObj := old.(*provision.ProvisionRequest)
	newObj.Spec = oldObj.Spec
	newObj.Labels = oldObj.Labels
}

func (provisionRequestStatusStrategy) ValidateUpdate(ctx context.Context, obj, old runtime.Object) field.ErrorList {
	return nil
}

func (provisionRequestStatusStrategy) WarningsOnUpdate(ctx context.Context, obj, old runtime.Object) []string {
	return nil
}
