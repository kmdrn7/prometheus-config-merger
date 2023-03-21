package config

import (
	"errors"
	"log"

	"github.com/spf13/viper"
)

type PrometheusConfig struct {
	Id         string            `mapstructure:"id"`
	Path       string            `mapstructure:"path"`
	Weight     int32             `mapstructure:"weight"`
	Kubernetes *KubernetesConfig `mapstructure:"kubernetes"`
}

type KubernetesConfig struct {
	Secret    *KubernetesConfigItem `mapstructure:"secret"`
	ConfigMap *KubernetesConfigItem `mapstructure:"configmap"`
}

type KubernetesConfigItem struct {
	Name      string `mapstructure:"name"`
	Key       string `mapstructure:"key"`
	Namespace string `mapstructure:"namespace"`
}

type Config struct {
	PrometheusEndpoint     string             `mapstructure:"prometheus_endpoint"`
	PrometheusConfigs      []PrometheusConfig `mapstructure:"prometheus_configs"`
	TargetPrometheusConfig string             `mapstructure:"target_prometheus_config"`
	ServerListenExternal   string             `mapstructure:"server_listen_external"`
	Debug                  bool
}

var localConfig *Config

func New() *Config {

	localConfig = &Config{
		PrometheusEndpoint:     "",
		PrometheusConfigs:      nil,
		TargetPrometheusConfig: "",
		Debug:                  false,
	}

	// run viper to parse config file
	if err := viper.Unmarshal(localConfig); err != nil {
		log.Fatal(err.Error())
	}

	return localConfig
}

func GetConfig() *Config {
	return localConfig
}

func GetServerListenExternal() string {
	return localConfig.ServerListenExternal
}

// Validate run viper to parse config file and inject to Config struct
func (config *Config) Validate() error {
	for _, prometheusConfig := range config.PrometheusConfigs {
		if prometheusConfig.Kubernetes != nil {
			if prometheusConfig.Kubernetes.ConfigMap != nil && prometheusConfig.Kubernetes.Secret != nil {
				return errors.New("cannot use kubernetes.configmap and kubernetes.secret at same time")
			}
		}
	}
	return nil
}
