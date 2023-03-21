package merge

import (
	"log"
	"os"
	"prometheus-config-merger/pkg/config"
	"prometheus-config-merger/pkg/utils"
	"sort"
	"sync"

	gokitlog "github.com/go-kit/log"
	prometheusconfig "github.com/prometheus/prometheus/config"

	// install custom prometheus service discovery
	_ "github.com/prometheus/prometheus/discovery/file"
	_ "github.com/prometheus/prometheus/discovery/gce"
	_ "github.com/prometheus/prometheus/discovery/kubernetes"
)

func RunWithReturn(cfg *config.Config) ([]byte, error) {
	content, err := run(cfg)
	if err != nil {
		panic(err)
	}
	return content, nil
}

func RunWithoutReturn(cfg *config.Config) {
	m := &sync.Mutex{}
	content, err := run(cfg)
	if err != nil {
		panic(err)
	}
	WriteToFileSafe(cfg.TargetPrometheusConfig, content, m)
}

func run(cfg *config.Config) ([]byte, error) {
	logger := gokitlog.NewNopLogger()

	sort.Slice(cfg.PrometheusConfigs, func(i, j int) bool {
		return cfg.PrometheusConfigs[i].Weight < cfg.PrometheusConfigs[j].Weight
	})

	finalPromeConfig := &prometheusconfig.Config{}

	for idx, promeConfig := range cfg.PrometheusConfigs {
		localPromeConfig, err := prometheusconfig.LoadFile(promeConfig.Path, false, false, logger)
		if err != nil {
			// ignore error for n+1 config file
			// send warning message for checking the correct prometheus config path
			if os.IsNotExist(err) && idx > 0 {
				log.Println("file", promeConfig.Path, "is not exist yet, skipping")
				log.Println("please make sure prometheus config path is correct!!!")
				continue
			}
			return []byte(""), err
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
		return []byte(""), err
	}

	return b.Bytes(), err
}

// WriteToFileSafe write content to filepath safely using golang mutex
func WriteToFileSafe(filepath string, content []byte, m *sync.Mutex) {
	// acquire mutex lock
	m.Lock()
	// open target file and rewrite the content
	f, err := os.Create(filepath)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer f.Close()
	if _, err := f.Write(content); err != nil {
		log.Fatal(err.Error())
	}
	// release mutex lock
	m.Unlock()
}
