package merge

import (
	"log"
	"os"
	"prometheus-config-merger/pkg/config"
	"prometheus-config-merger/pkg/utils"
	"sort"

	prometheusconfig "github.com/prometheus/prometheus/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

// rootCmd represents the base command when called without any subcommands
var MergeCmd = &cobra.Command{
	Use:   "merge",
	Short: "Merge multiple prometheus config files into single config file.",
	Long: `
Merge rule_files and scrape_configs across multiple prometheus config files into single config file.
Just in case you want to segregate prometheus configuration into multiple file, and have it merged before consumed by prometheus-server.
`,
	Run: func(cmd *cobra.Command, args []string) {

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
	},
}
