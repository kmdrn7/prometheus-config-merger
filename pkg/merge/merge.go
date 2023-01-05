package merge

import (
	"log"
	"os"
	"prometheus-config-merger/pkg/config"
	"prometheus-config-merger/pkg/utils"
	"sort"

	prometheusconfig "github.com/prometheus/prometheus/config"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

func Run() {
	cfg := &config.Config{}
	if err := viper.Unmarshal(cfg); err != nil {
		log.Fatal(err.Error())
	}

	sort.Slice(cfg.PrometheusConfigs, func(i, j int) bool {
		return cfg.PrometheusConfigs[i].Weight < cfg.PrometheusConfigs[j].Weight
	})

	finalPromeConfig := &prometheusconfig.Config{}

	for idx, promeConfig := range cfg.PrometheusConfigs {
		localPromeConfig := &prometheusconfig.Config{}
		yfile, err := os.ReadFile(promeConfig.Path)
		if err != nil {
			// ignore error for n+1 config file
			// send warning message for checking the correct prometheus config path
			if os.IsNotExist(err) && idx > 0 {
				log.Println("file", promeConfig.Path, "is not exist yet, skipping")
				log.Println("please make sure prometheus config path is correct!!!")
				continue
			}
			log.Fatal(err.Error())
		}
		if err := yaml.Unmarshal(yfile, &localPromeConfig); err != nil {
			log.Fatal(err.Error())
		}

		if idx == 0 {
			finalPromeConfig = localPromeConfig
		} else {
			finalPromeConfig.RuleFiles = append(finalPromeConfig.RuleFiles, localPromeConfig.RuleFiles...)
			finalPromeConfig.ScrapeConfigs = append(finalPromeConfig.ScrapeConfigs, localPromeConfig.ScrapeConfigs...)
		}
	}

	b, err := utils.YamlEncode(&finalPromeConfig)
	if err != nil {
		log.Fatal(err.Error())
	}

	// open target file and rewrite the content
	f, err := os.Create(cfg.TargetPrometheusConfig)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer f.Close()
	if _, err := f.Write(b.Bytes()); err != nil {
		log.Fatal(err.Error())
	}
}
