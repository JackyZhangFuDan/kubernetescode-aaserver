/*
Copyright The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by lister-gen. DO NOT EDIT.

package v1alpha1

import (
	v1alpha1 "github.com/kubernetescode-aaserver/pkg/apis/provision/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// ProvisionRequestLister helps list ProvisionRequests.
// All objects returned here must be treated as read-only.
type ProvisionRequestLister interface {
	// List lists all ProvisionRequests in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha1.ProvisionRequest, err error)
	// ProvisionRequests returns an object that can list and get ProvisionRequests.
	ProvisionRequests(namespace string) ProvisionRequestNamespaceLister
	ProvisionRequestListerExpansion
}

// provisionRequestLister implements the ProvisionRequestLister interface.
type provisionRequestLister struct {
	indexer cache.Indexer
}

// NewProvisionRequestLister returns a new ProvisionRequestLister.
func NewProvisionRequestLister(indexer cache.Indexer) ProvisionRequestLister {
	return &provisionRequestLister{indexer: indexer}
}

// List lists all ProvisionRequests in the indexer.
func (s *provisionRequestLister) List(selector labels.Selector) (ret []*v1alpha1.ProvisionRequest, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.ProvisionRequest))
	})
	return ret, err
}

// ProvisionRequests returns an object that can list and get ProvisionRequests.
func (s *provisionRequestLister) ProvisionRequests(namespace string) ProvisionRequestNamespaceLister {
	return provisionRequestNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// ProvisionRequestNamespaceLister helps list and get ProvisionRequests.
// All objects returned here must be treated as read-only.
type ProvisionRequestNamespaceLister interface {
	// List lists all ProvisionRequests in the indexer for a given namespace.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha1.ProvisionRequest, err error)
	// Get retrieves the ProvisionRequest from the indexer for a given namespace and name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1alpha1.ProvisionRequest, error)
	ProvisionRequestNamespaceListerExpansion
}

// provisionRequestNamespaceLister implements the ProvisionRequestNamespaceLister
// interface.
type provisionRequestNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all ProvisionRequests in the indexer for a given namespace.
func (s provisionRequestNamespaceLister) List(selector labels.Selector) (ret []*v1alpha1.ProvisionRequest, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.ProvisionRequest))
	})
	return ret, err
}

// Get retrieves the ProvisionRequest from the indexer for a given namespace and name.
func (s provisionRequestNamespaceLister) Get(name string) (*v1alpha1.ProvisionRequest, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1alpha1.Resource("provisionrequest"), name)
	}
	return obj.(*v1alpha1.ProvisionRequest), nil
}