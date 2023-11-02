package controller

import (
	"context"
	"fmt"
	"time"

	"github.com/kubernetescode-aaserver/pkg/apis/provision/v1alpha1"
	prclientset "github.com/kubernetescode-aaserver/pkg/generated/clientset/versioned"
	prInformer "github.com/kubernetescode-aaserver/pkg/generated/informers/externalversions/provision/v1alpha1"
	prlist "github.com/kubernetescode-aaserver/pkg/generated/listers/provision/v1alpha1"

	apps "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/uuid"
	"k8s.io/apimachinery/pkg/util/wait"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog/v2"
)

type Controller struct {
	coreAPIClient clientset.Interface
	prClient      prclientset.Interface

	prLister prlist.ProvisionRequestLister
	prSynced cache.InformerSynced

	queue       workqueue.RateLimitingInterface
	syncHandler func(ctx context.Context, key string) error
}

func NewProvisionController(prClient prclientset.Interface, prInfo prInformer.ProvisionRequestInformer,
	coreAPIClient clientset.Interface) *Controller {

	c := &Controller{
		coreAPIClient: coreAPIClient,
		prClient:      prClient,
		prLister:      prInfo.Lister(),
		prSynced:      prInfo.Informer().HasSynced,

		queue: workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "provisionrequest"),
	}
	c.syncHandler = c.sync

	prInfo.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		// when create
		AddFunc: func(obj interface{}) {
			klog.Info("New Provision Request is found")
			cast := obj.(*v1alpha1.ProvisionRequest)
			key, err := cache.MetaNamespaceKeyFunc(cast)
			if err != nil {
				klog.ErrorS(err, "Failed when extracting key of Provision Request Object")
				return
			}
			c.queue.Add(key)
		},
	})
	return c
}

func (c *Controller) Run(ctx context.Context, workers int) {
	defer utilruntime.HandleCrash()
	defer c.queue.ShutDown()

	klog.Info("Starting Provision Controller")
	defer klog.Info("Shutting down Provision Controller")

	klog.Info("Waiting for caches to sync for Provision controller")
	if !cache.WaitForCacheSync(ctx.Done(), c.prSynced) {
		utilruntime.HandleError(fmt.Errorf("unable to sync caches for provision controller"))
		return
	}
	klog.Infof("Caches are synced for Provision controller")

	// wait.Until的作用是如果runWorker处理失败退出，那么再次启动它
	for i := 0; i < workers; i++ {
		go wait.UntilWithContext(ctx, c.runWorker, time.Second)
	}

	<-ctx.Done()
}

func (c *Controller) runWorker(ctx context.Context) {
	for c.processNextWorkItem(ctx) {
	}
}

func (c *Controller) processNextWorkItem(ctx context.Context) bool {
	key, quit := c.queue.Get()
	if quit {
		return false
	}
	defer c.queue.Done(key)

	err := c.syncHandler(ctx, key.(string))
	if err == nil {
		c.queue.Forget(key)
		klog.Infof("Finish processing key %s", key)
		return true
	}
	utilruntime.HandleError(fmt.Errorf("%v failed with: %v", key, err))
	c.queue.AddRateLimited(key)
	return true
}

func (c *Controller) sync(ctx context.Context, key string) (err error) {
	klog.Infof("start to run sync logic for PR %s", key)
	defer klog.Infof("finish sync logic for PR %s", key)

	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		klog.ErrorS(err, "Failed to split meta namespace cache key", "cacheKey", name)
	}

	pr, err := c.prLister.ProvisionRequests(namespace).Get(name)
	if errors.IsNotFound(err) {
		klog.Infof("Provision Request %s has been deleted", key)
		return nil
	}
	if err != nil {
		return err
	}

	pr2 := pr.DeepCopy()

	custNameSpaceName := pr.Spec.NamespaceName
	_, err = c.coreAPIClient.CoreV1().Namespaces().Get(ctx, custNameSpaceName, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		custNameSpace := v1.Namespace{
			TypeMeta: metav1.TypeMeta{APIVersion: "v1", Kind: "Namespace"},
			ObjectMeta: metav1.ObjectMeta{
				UID:         uuid.NewUUID(),
				Name:        custNameSpaceName,
				Annotations: make(map[string]string),
			},
			Spec: v1.NamespaceSpec{},
		}
		_, err = c.coreAPIClient.CoreV1().Namespaces().Create(ctx, &custNameSpace, metav1.CreateOptions{})
		if err != nil {
			return &errors.StatusError{ErrStatus: metav1.Status{
				Status:  "Failure",
				Message: "fail to create customer namespace",
			}}
		}
	}

	if !pr.Status.DbReady {
		var replicas int32 = 1
		selector := map[string]string{}
		selector["type"] = "provisioinrequest"
		selector["company"] = pr.Labels["company"]

		d := apps.Deployment{
			TypeMeta: metav1.TypeMeta{APIVersion: "apps/v1", Kind: "Deployment"},
			ObjectMeta: metav1.ObjectMeta{
				UID:         uuid.NewUUID(),
				Name:        "cust-db",
				Namespace:   custNameSpaceName,
				Annotations: make(map[string]string),
			},
			Spec: apps.DeploymentSpec{
				Replicas: &replicas,
				Selector: &metav1.LabelSelector{MatchLabels: selector},
				Template: v1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: selector,
					},
					Spec: v1.PodSpec{
						Containers: []v1.Container{
							{
								Name:            "customer-db",
								Image:           "mysql:5.7",
								ImagePullPolicy: "IfNotPresent",
								Env:             []v1.EnvVar{{Name: "MYSQL_ROOT_PASSWORD", Value: "pleasechangetosecret"}},
								Ports:           []v1.ContainerPort{{ContainerPort: 3306, Name: "mysql"}},
							},
						},
					},
				},
			},
		}
		_, err = c.coreAPIClient.AppsV1().Deployments(custNameSpaceName).Get(ctx, d.Name, metav1.GetOptions{})
		if errors.IsNotFound(err) {
			_, err = c.coreAPIClient.AppsV1().Deployments(custNameSpaceName).Create(ctx, &d, metav1.CreateOptions{})
			if err != nil {
				klog.ErrorS(err, "Failed when creating DB deployment for Provision Request")
				return err
			}
		} else if err != nil {
			return &errors.StatusError{ErrStatus: metav1.Status{
				Status:  "Failure",
				Message: "fail to read DB deployment",
			}}
		}
	}

	if !pr.Status.IngressReady {
		// 这里省去配置Ingress的逻辑......
	}

	pr2.Status.IngressReady = true
	pr2.Status.DbReady = true
	pr2.Kind = "ProvisionRequest"
	_, err = c.prClient.ProvisionV1alpha1().ProvisionRequests(pr2.Namespace).UpdateStatus(
		context.TODO(), pr2, metav1.UpdateOptions{})
	if err != nil {
		klog.ErrorS(err, "Fail to update request status")
		return &errors.StatusError{ErrStatus: metav1.Status{
			Status:  "Failure",
			Message: "fail to update provision request status",
		}}
	}

	klog.Infof("Sucessfully fulfill provision request %s", key)
	return nil
}
