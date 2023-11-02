package plugin

import (
	"context"
	"fmt"
	"io"

	"github.com/kubernetescode-aaserver/pkg/apis/provision"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
	"k8s.io/apiserver/pkg/admission"
	"k8s.io/klog/v2"

	informers "github.com/kubernetescode-aaserver/pkg/generated/informers/externalversions"
	listers "github.com/kubernetescode-aaserver/pkg/generated/listers/provision/v1alpha1"
)

type ProvisionPlugin struct {
	*admission.Handler
	Lister listers.ProvisionRequestLister
}

// method defined by admission.ValidationInterface
func (plugin *ProvisionPlugin) Validate(ctx context.Context, a admission.Attributes,
	interfaces admission.ObjectInterfaces) error {
	klog.Info("provision admission plugin's validate method starts")

	if a.GetOperation() != admission.Create {
		return nil
	}

	if a.GetKind().GroupKind() != provision.Kind("ProvisionRequest") {
		return nil
	}

	if !plugin.WaitForReady() {
		return admission.NewForbidden(a,
			fmt.Errorf("the plugin isn't ready for handling request"))
	}

	metaAccessor, err := meta.Accessor(a.GetObject())
	if err != nil {
		return err
	}
	company := metaAccessor.GetLabels()["company"]
	req, err := labels.NewRequirement("company", selection.Equals, []string{company})
	if err != nil {
		return admission.NewForbidden(a,
			fmt.Errorf("failed to create label requirement"))
	}

	reqs, err := plugin.Lister.List(labels.NewSelector().Add(*req))
	if len(reqs) > 0 {
		return admission.NewForbidden(a,
			fmt.Errorf("the company already has provision request"))
	}
	return nil
}

func New() (*ProvisionPlugin, error) {
	return &ProvisionPlugin{
		Handler: admission.NewHandler(admission.Create),
	}, nil
}

func Register(plugin *admission.Plugins) {
	plugin.Register("Provision", func(config io.Reader) (admission.Interface, error) {
		return New()
	})
}

func (plugin *ProvisionPlugin) SetInformerFactory(f informers.SharedInformerFactory) {
	plugin.Lister = f.Provision().V1alpha1().ProvisionRequests().Lister()
	plugin.SetReadyFunc(f.Provision().V1alpha1().ProvisionRequests().Informer().HasSynced)
}
