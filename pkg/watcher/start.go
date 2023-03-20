package watcher

import (
	"fmt"
	"log"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

func (w *Watcher) Run() error {
	w.setupConfig()
	w.setupWatchedResources()
	w.setupKubernetesClientset()
	w.StartWatcher()
	return nil
}

func (w *Watcher) setupKubernetesClientset() {
	// check for available kubenetes configuration
	config, err := clientcmd.BuildConfigFromFlags("", w.kubeconfig)
	if err != nil {
		log.Panic(err)
	}
	// build kubernetes clientset from avaiable config
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}
	w.clientset = clientset
}

func (w *Watcher) StartWatcher() {
	// stop signal for the informer
	stopper := make(chan struct{})
	defer close(stopper)

	// setup shared informers
	factory := informers.NewSharedInformerFactoryWithOptions(w.clientset, 0)
	secretInformer := factory.Core().V1().Secrets().Informer()
	configmapInformer := factory.Core().V1().ConfigMaps().Informer()

	// handle runtime crash
	defer runtime.HandleCrash()

	// start informer ->
	go factory.Start(stopper)

	// start to sync secretInformer and call list
	if !cache.WaitForCacheSync(stopper, secretInformer.HasSynced) {
		runtime.HandleError(fmt.Errorf("timed out waiting for secretInformer caches to sync"))
		return
	}

	// start to sync configmapInformer and call list
	if !cache.WaitForCacheSync(stopper, configmapInformer.HasSynced) {
		runtime.HandleError(fmt.Errorf("timed out waiting for configmapInformer caches to sync"))
		return
	}

	// add handler for secret event
	secretInformer.AddEventHandler(
		cache.FilteringResourceEventHandler{
			Handler: cache.ResourceEventHandlerFuncs{
				AddFunc:    secretOnAdd,
				UpdateFunc: secretOnUpdate,
				DeleteFunc: secretOnDelete,
			},
			FilterFunc: func(obj interface{}) bool {
				secret, ok := obj.(*corev1.Secret)
				if !ok {
					return true
				}
				return isAllowedResource("Secret", secret.Name, secret.Namespace)
			},
		},
	)

	// add handler for configmap event
	configmapInformer.AddEventHandler(
		cache.FilteringResourceEventHandler{
			Handler: cache.ResourceEventHandlerFuncs{
				AddFunc:    configmapOnAdd,
				UpdateFunc: configmapOnUpdate,
				DeleteFunc: configmapOnDelete,
			},
			FilterFunc: func(obj interface{}) bool {
				configmap, ok := obj.(*corev1.ConfigMap)
				if !ok {
					return true
				}
				return isAllowedResource("ConfigMap", configmap.Name, configmap.Namespace)
			},
		},
	)

	// handle stop signal
	<-stopper
}

// isAllowedResource check secret/confirmap agains watchedResourceList
func isAllowedResource(resourcetype string, resourcename string, resourcenamespace string) bool {
	for _, wr := range watchedResourceList {
		// only filter for resource type
		if wr.Type != resourcetype {
			continue
		}
		// which match the resource name
		if resourcename != wr.Name {
			continue
		}
		// which match the resource namespace
		if resourcenamespace != wr.Namespace {
			continue
		}
		return true
	}
	return false
}
