package apiserver

import (
	"os"
	"path/filepath"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/version"
	openapinamer "k8s.io/apiserver/pkg/endpoints/openapi"
	"k8s.io/apiserver/pkg/registry/rest"
	gserver "k8s.io/apiserver/pkg/server"
	"k8s.io/apiserver/pkg/util/openapi"
	"k8s.io/client-go/kubernetes"
	clientgorest "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"

	"github.com/kubernetescode-aaserver/pkg/apis/provision"
	"github.com/kubernetescode-aaserver/pkg/apis/provision/install"
	prcontroller "github.com/kubernetescode-aaserver/pkg/controller"
	prclientset "github.com/kubernetescode-aaserver/pkg/generated/clientset/versioned"
	prinformerfactory "github.com/kubernetescode-aaserver/pkg/generated/informers/externalversions"
	generatedopenapi "github.com/kubernetescode-aaserver/pkg/generated/openapi"
	"github.com/kubernetescode-aaserver/pkg/registry"
	provisionstore "github.com/kubernetescode-aaserver/pkg/registry/provision"
)

var (
	Scheme = runtime.NewScheme()
	Codecs = serializer.NewCodecFactory(Scheme)
)

func init() {
	install.Install(Scheme)
	metav1.AddToGroupVersion(Scheme, schema.GroupVersion{Version: "v1"})
	unversioned := schema.GroupVersion{Group: "", Version: "v1"}
	Scheme.AddUnversionedTypes(
		unversioned,
		&metav1.Status{},
		&metav1.APIVersions{},
		&metav1.APIGroupList{},
		&metav1.APIGroup{},
		&metav1.APIResourceList{},
	)
}

type MyServer struct {
	GenericAPIServer *gserver.GenericAPIServer
}

type Config struct {
	GenericConfig *gserver.RecommendedConfig
}

type completedConfig struct {
	GenericConfig gserver.CompletedConfig
}

type CompletedConfig struct {
	*completedConfig
}

func (cfg *Config) Complete() CompletedConfig {
	// version
	cconfig := completedConfig{
		cfg.GenericConfig.Complete(),
	}
	cconfig.GenericConfig.Version = &version.Info{
		Major: "1",
		Minor: "0",
	}
	// openapiv3
	getOpenAPIDefinitions := openapi.GetOpenAPIDefinitionsWithoutDisabledFeatures(generatedopenapi.GetOpenAPIDefinitions)
	cconfig.GenericConfig.OpenAPIV3Config = gserver.DefaultOpenAPIV3Config(getOpenAPIDefinitions, openapinamer.NewDefinitionNamer(Scheme))
	cconfig.GenericConfig.OpenAPIV3Config.Info.Title = "aggregated-apiserver"

	return CompletedConfig{&cconfig}
}

func (ccfg completedConfig) NewServer() (*MyServer, error) {
	genericServer, err := ccfg.GenericConfig.New(
		"provision-apiserver",
		gserver.NewEmptyDelegate())
	if err != nil {
		return nil, err
	}

	server := &MyServer{
		GenericAPIServer: genericServer,
	}

	apiGroupInfo := gserver.NewDefaultAPIGroupInfo(
		provision.GroupName,
		Scheme,
		metav1.ParameterCodec,
		Codecs,
	)
	v1alphastorage := map[string]rest.Storage{}
	prRest, prStatusRest, err := provisionstore.NewREST(Scheme, ccfg.GenericConfig.RESTOptionsGetter)
	v1alphastorage["provisionrequests"] = registry.RESTWithErrorHandler(prRest, err)
	v1alphastorage["provisionrequests/status"] = registry.RESTWithErrorHandler(prStatusRest, err)
	apiGroupInfo.VersionedResourcesStorageMap["v1alpha1"] = v1alphastorage

	if err := server.GenericAPIServer.InstallAPIGroup(&apiGroupInfo); err != nil {
		return nil, err
	}

	// controller
	config, err := clientgorest.InClusterConfig()
	if err != nil {
		// fallback to kubeconfig
		kubeconfig := filepath.Join("~", ".kube", "config")
		if envvar := os.Getenv("KUBECONFIG"); len(envvar) > 0 {
			kubeconfig = envvar
		}
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			klog.ErrorS(err, "The kubeconfig cannot be loaded: %v\n")
			panic(err)
		}
	}
	coreAPIClientset, err := kubernetes.NewForConfig(config)

	client, err := prclientset.NewForConfig(genericServer.LoopbackClientConfig)
	if err != nil {
		klog.Error("Can't create client set for provision during creating server")
	}
	prInformerFactory := prinformerfactory.NewSharedInformerFactory(client, 0)
	controller := prcontroller.NewProvisionController(
		client,
		prInformerFactory.Provision().V1alpha1().ProvisionRequests(),
		coreAPIClientset)

	genericServer.AddPostStartHookOrDie("aapiserver-controller", func(ctx gserver.PostStartHookContext) error {
		ctxpr := wait.ContextForChannel(ctx.StopCh)
		go func() {
			controller.Run(ctxpr, 2)
		}()
		return nil
	})
	genericServer.AddPostStartHookOrDie("aapiserver-informer", func(context gserver.PostStartHookContext) error {
		prInformerFactory.Start(context.StopCh)
		return nil
	})
	return server, nil
}
