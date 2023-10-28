package apiserver

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apimachinery/pkg/version"
	"k8s.io/apiserver/pkg/registry/rest"
	gserver "k8s.io/apiserver/pkg/server"

	"github.com/kubernetescode-aaserver/pkg/apis/provision"
	"github.com/kubernetescode-aaserver/pkg/apis/provision/install"
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
	cconfig := completedConfig{
		cfg.GenericConfig.Complete(),
	}
	cconfig.GenericConfig.Version = &version.Info{
		Major: "1",
		Minor: "0",
	}
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
	v1alphastorage["provisionrequests"] = registry.RESTWithErrorHandler(
		provisionstore.NewREST(Scheme, ccfg.GenericConfig.RESTOptionsGetter))
	apiGroupInfo.VersionedResourcesStorageMap["v1alpha1"] = v1alphastorage

	if err := server.GenericAPIServer.InstallAPIGroup(&apiGroupInfo); err != nil {
		return nil, err
	}

	return server, nil
}
