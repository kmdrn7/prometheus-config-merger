package config

import (
	"log"

	"github.com/spf13/viper"
)

type PrometheusConfig struct {
	Id     string `mapstructure:"id"`
	Path   string `mapstructure:"path"`
	Weight int32  `mapstructure:"weight"`
}

type Config struct {
	PrometheusEndpoint     string             `mapstructure:"prometheus_endpoint"`
	PrometheusConfigs      []PrometheusConfig `mapstructure:"prometheus_configs"`
	TargetPrometheusConfig string             `mapstructure:"target_prometheus_config"`
	Debug                  bool
}

func GetConfig() *Config {
	cfg := &Config{}
	if err := viper.Unmarshal(cfg); err != nil {
		log.Fatal(err.Error())
	}
	return cfg
}
