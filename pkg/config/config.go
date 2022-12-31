package config

type PrometheusConfig struct {
	Id     string `mapstructure:"id"`
	Path   string `mapstructure:"path"`
	Weight int32  `mapstructure:"weight"`
}

type Config struct {
	PrometheusConfigs      []PrometheusConfig `mapstructure:"prometheus_configs"`
	TargetPrometheusConfig string             `mapstructure:"target_prometheus_config"`
}
