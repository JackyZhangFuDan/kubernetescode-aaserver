package cmd

import (
	"fmt"
	"net"

	myadmission "github.com/kubernetescode-aaserver/pkg/admission"
	"github.com/kubernetescode-aaserver/pkg/admission/plugin"
	v1alpha1 "github.com/kubernetescode-aaserver/pkg/apis/provision/v1alpha1"
	apiserver "github.com/kubernetescode-aaserver/pkg/apiserver"
	informers "github.com/kubernetescode-aaserver/pkg/generated/informers/externalversions"

	clientset "github.com/kubernetescode-aaserver/pkg/generated/clientset/versioned"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apiserver/pkg/admission"
	gserver "k8s.io/apiserver/pkg/server"
	genericoptions "k8s.io/apiserver/pkg/server/options"
)

type ServerOptions struct {
	RecommendedOptions    *genericoptions.RecommendedOptions
	SharedInformerFactory informers.SharedInformerFactory
}

func NewServerOptions() *ServerOptions {
	o := &ServerOptions{
		RecommendedOptions: genericoptions.NewRecommendedOptions(
			"/registry/provision-apiserver.mydomain.com",
			apiserver.Codecs.LegacyCodec(v1alpha1.SchemeGroupVersion),
		),
	}
	return o
}

func (o *ServerOptions) Complete() error {
	// adminition plugin
	plugin.Register(o.RecommendedOptions.Admission.Plugins)
	o.RecommendedOptions.Admission.RecommendedPluginOrder =
		append(o.RecommendedOptions.Admission.RecommendedPluginOrder, "Provision")
	// admission plugin initializer
	o.RecommendedOptions.ExtraAdmissionInitializers =
		func(cfg *gserver.RecommendedConfig) ([]admission.PluginInitializer, error) {
			client, err := clientset.NewForConfig(cfg.LoopbackClientConfig)
			if err != nil {
				return nil, err
			}
			informerFactory := informers.NewSharedInformerFactory(client, cfg.LoopbackClientConfig.Timeout)
			o.SharedInformerFactory = informerFactory
			return []admission.PluginInitializer{myadmission.NewProvisionPluginInitializer(informerFactory)},
				nil
		}
	return nil
}

func (o *ServerOptions) Validate() error {
	errors := []error{}
	errors = append(errors, o.RecommendedOptions.Validate()...)
	return utilerrors.NewAggregate(errors)
}

func (o *ServerOptions) Config() (*apiserver.Config, error) {
	if err := o.RecommendedOptions.SecureServing.MaybeDefaultWithSelfSignedCerts(
		"localhost", nil, []net.IP{net.ParseIP("127.0.0.1")}); err != nil {
		return nil, fmt.Errorf("error creating self-signed certificates: %v", err)
	}

	standardConfig := gserver.NewRecommendedConfig(apiserver.Codecs)
	if err := o.RecommendedOptions.ApplyTo(standardConfig); err != nil {
		return nil, err
	}
	myConfig := &apiserver.Config{
		GenericConfig: standardConfig,
	}
	return myConfig, nil
}
