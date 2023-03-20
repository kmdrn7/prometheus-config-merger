package watcher

import (
	"prometheus-config-merger/pkg/config"

	"golang.org/x/net/context"
	"k8s.io/client-go/kubernetes"
)

type Watcher struct {
	incluster  bool
	kubeconfig string

	clientset kubernetes.Interface
	ctx       context.Context
	cfg       *config.Config
}

type WatchedResource struct {
	Type      string
	Name      string
	Namespace string
	Key       string
	Path      string
}

var watchedResourceList []WatchedResource

func NewWatcher(kubeconfig string) *Watcher {
	return &Watcher{
		clientset:  nil,
		incluster:  true,
		kubeconfig: expandKubeconfig(kubeconfig),
		ctx:        context.TODO(),
	}
}
