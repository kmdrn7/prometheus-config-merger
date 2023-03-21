package watcher

import (
	"errors"
	"fmt"
	"log"
	"prometheus-config-merger/pkg/config"
	"prometheus-config-merger/pkg/http"
	"prometheus-config-merger/pkg/merge"
	"prometheus-config-merger/pkg/utils"

	corev1 "k8s.io/api/core/v1"
)

func secretOnAdd(obj interface{}) {
	// acquire mutex lock
	mut.Lock()

	secret := obj.(*corev1.Secret)
	wr, err := getWatchedResourceMapping("Secret", secret.Name, secret.Namespace)
	if err != nil {
		panic(err)
	}

	log.Printf("event [secretOnAdd] triggered for [%s/%s] \n", wr.Namespace, wr.Name)
	log.Printf("syncing secret [%s/%s] to %s \n", wr.Namespace, wr.Name, wr.Path)
	if err := utils.SyncResourceContentToLocalFile(secret.Data[wr.Key], wr.Path); err != nil {
		panic(err)
	}

	cfg := config.GetConfig()
	content, err := merge.RunWithReturn(cfg)
	if err != nil {
		panic(err)
	}
	if err := syncToExternalService(string(content), cfg.TargetPrometheusConfig); err != nil {
		panic(err)
	}

	// release mutex lock
	mut.Unlock()
}

func secretOnUpdate(oldObj interface{}, newObj interface{}) {
	// acquire mutex lock
	mut.Lock()

	secret := newObj.(*corev1.Secret)
	wr, err := getWatchedResourceMapping("Secret", secret.Name, secret.Namespace)
	if err != nil {
		panic(err)
	}
	log.Printf("event [secretOnUpdate] triggered for [%s/%s] \n", wr.Namespace, wr.Name)
	log.Printf("syncing secret [%s/%s] to %s \n", wr.Namespace, wr.Name, wr.Path)
	if err := utils.SyncResourceContentToLocalFile(secret.Data[wr.Key], wr.Path); err != nil {
		panic(err)
	}

	cfg := config.GetConfig()
	content, err := merge.RunWithReturn(cfg)
	if err != nil {
		panic(err)
	}
	if err := syncToExternalService(string(content), cfg.TargetPrometheusConfig); err != nil {
		panic(err)
	}

	// release mutex lock
	mut.Unlock()
}

func secretOnDelete(obj interface{}) {
	secret := obj.(*corev1.Secret)
	log.Printf("event [secretOnDelete] triggered for [%s/%s] \n", secret.Namespace, secret.Name)
}

func configmapOnAdd(obj interface{}) {
	// acquire mutex lock
	mut.Lock()

	configmap := obj.(*corev1.ConfigMap)
	wr, err := getWatchedResourceMapping("ConfigMap", configmap.Name, configmap.Namespace)
	if err != nil {
		panic(err)
	}
	log.Printf("event [configmapOnAdd] triggered for [%s/%s] \n", wr.Namespace, wr.Name)
	log.Printf("syncing configmap [%s/%s] to %s \n", wr.Namespace, wr.Name, wr.Path)
	if err := utils.SyncResourceContentToLocalFile([]byte(configmap.Data[wr.Key]), wr.Path); err != nil {
		panic(err)
	}

	cfg := config.GetConfig()
	content, err := merge.RunWithReturn(cfg)
	if err != nil {
		panic(err)
	}
	if err := syncToExternalService(string(content), cfg.TargetPrometheusConfig); err != nil {
		panic(err)
	}

	// release mutex lock
	mut.Unlock()
}

func configmapOnUpdate(oldObj interface{}, newObj interface{}) {
	// acquire mutex lock
	mut.Lock()

	configmap := newObj.(*corev1.ConfigMap)
	wr, err := getWatchedResourceMapping("ConfigMap", configmap.Name, configmap.Namespace)
	if err != nil {
		panic(err)
	}
	log.Printf("event [configmapOnUpdate] triggered for [%s/%s] \n", wr.Namespace, wr.Name)
	log.Printf("syncing configmap [%s/%s] to %s \n", wr.Namespace, wr.Name, wr.Path)
	if err := utils.SyncResourceContentToLocalFile([]byte(configmap.Data[wr.Key]), wr.Path); err != nil {
		panic(err)
	}

	cfg := config.GetConfig()
	content, err := merge.RunWithReturn(cfg)
	if err != nil {
		panic(err)
	}
	if err := syncToExternalService(string(content), cfg.TargetPrometheusConfig); err != nil {
		panic(err)
	}

	// release mutex lock
	mut.Unlock()
}

func configmapOnDelete(obj interface{}) {
	configmap := obj.(*corev1.ConfigMap)
	log.Printf("event [configmapOnDelete] triggered for [%s/%s] \n", configmap.Namespace, configmap.Name)
}

// getWatchedResourceMapping get mapping of secret/configmap from watchedresource list
func getWatchedResourceMapping(resourcetype string, secretName string, secretNamespace string) (WatchedResource, error) {
	for _, wr := range watchedResourceList {
		// filter for resource type, name, and namespace
		if wr.Type == resourcetype && secretName == wr.Name && secretNamespace == wr.Namespace {
			return wr, nil
		}
	}
	return WatchedResource{}, fmt.Errorf("requested %s is not listed in watchedresource list", resourcetype)
}

// syncToExternalService send final prometheus config payload to watcher
func syncToExternalService(content string, filepath string) error {
	res, err := http.PostWithBody(externalService, ResourceHTTPBodyPayload{
		Content:  content,
		Filepath: filepath,
	})
	if err != nil {
		return err
	}
	if res.StatusCode() != http.StatusOK {
		return errors.New("failed sync payload to externalService")
	}
	return nil
}
