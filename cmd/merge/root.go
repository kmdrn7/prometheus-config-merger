package merge

import (
	"prometheus-config-merger/pkg/config"
	"prometheus-config-merger/pkg/merge"

	"github.com/spf13/cobra"
)

// Cmd represents the base command when called without any subcommands
var Cmd = &cobra.Command{
	Use:   "merge",
	Short: "Merge multiple prometheus config files into single config file.",
	Long: `
Merge rule_files and scrape_configs across multiple prometheus config files into single config file.
Just in case you want to segregate prometheus configuration into multiple file, and have it merged before consumed by prometheus-server.
`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.GetConfig()
		merge.RunWithReturn(cfg)
	},
}
