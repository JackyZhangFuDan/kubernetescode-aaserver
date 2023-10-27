package admission

import (
	informers "github.com/kubernetescode-aaserver/pkg/generated/informers/externalversions"
	"k8s.io/apiserver/pkg/admission"
)

type WantsInformerFactory interface {
	SetInformerFactory(informers.SharedInformerFactory)
}

type provisionPluginInitializer struct {
	informerFactory informers.SharedInformerFactory
}

func (i provisionPluginInitializer) Initialize(plugin admission.Interface) {
	if wants, ok := plugin.(WantsInformerFactory); ok {
		wants.SetInformerFactory(i.informerFactory)
	}
}

func NewProvisionPluginInitializer(
	informerFactory informers.SharedInformerFactory) provisionPluginInitializer {

	return provisionPluginInitializer{
		informerFactory: informerFactory,
	}
}
