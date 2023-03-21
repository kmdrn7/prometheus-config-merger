package watcher

import (
	"log"
	"os"
	"prometheus-config-merger/pkg/config"
)

const KUBECONFIG = "KUBECONFIG"

// expandKubeconfig try locate kubeconfig from environment variable when kubeconfig opt is empty
func expandKubeconfig(kubeconfig string) string {
	if kubeconfig == "" {
		if kubeconfigenv, exists := os.LookupEnv(KUBECONFIG); !exists {
			return ""
		} else {
			return kubeconfigenv
		}
	}
	return kubeconfig
}

func (w *Watcher) setupWatchedResources() {
	log.Println("setup watcher for kubernetes configmap/secret resources ...")
	for _, prometheusConfig := range w.cfg.PrometheusConfigs {
		if prometheusConfig.Kubernetes != nil {
			if prometheusConfig.Kubernetes.Secret != nil {
				watchedResourceList = append(watchedResourceList, WatchedResource{
					Type:      "Secret",
					Namespace: prometheusConfig.Kubernetes.Secret.Namespace,
					Name:      prometheusConfig.Kubernetes.Secret.Name,
					Key:       prometheusConfig.Kubernetes.Secret.Key,
					Path:      prometheusConfig.Path,
				})
			} else if prometheusConfig.Kubernetes.ConfigMap != nil {
				watchedResourceList = append(watchedResourceList, WatchedResource{
					Type:      "ConfigMap",
					Namespace: prometheusConfig.Kubernetes.ConfigMap.Namespace,
					Name:      prometheusConfig.Kubernetes.ConfigMap.Name,
					Key:       prometheusConfig.Kubernetes.ConfigMap.Key,
					Path:      prometheusConfig.Path,
				})
			}
		}
	}

	log.Println("start watcher for kubernetes configmap/secret resources ...")
	for _, resource := range watchedResourceList {
		log.Printf("... watching %s ...\n", resource)
	}
}

func (w *Watcher) setupConfig() {
	w.cfg = config.GetConfig()

	// configure external endpoint service
	externalService = w.cfg.ServerListenExternal + "/-/inject-config"
}
