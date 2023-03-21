package watcher

import (
	"prometheus-config-merger/pkg/watcher"

	"github.com/spf13/cobra"
)

var (
	dryRun     bool
	kubeconfig string
)

var w = watcher.NewWatcher(kubeconfig)

// Cmd represents the base command when called without any subcommands
var Cmd = &cobra.Command{
	Use:   "watcher",
	Short: "Watch configmap/secret containing prometheus config.",
	Long: `
Watcher will watch kubernetes secret/configmap changes that containing prometheus configuration
whether in plain text format or in gzipped format, then it will map and sync with the local file inside prometheus pod.
`,
	Run: func(cmd *cobra.Command, args []string) {
		if dryRun {
			cobra.CheckErr(w.DryRun())
		} else {
			cobra.CheckErr(w.Start())
		}
	},
}

func init() {
	Cmd.Flags().BoolVar(&dryRun, "dryrun", false, "Only print watcher configuration information. Usefull for debugging")
	Cmd.Flags().StringVar(&kubeconfig, "kubeconfig", "", "Specify where the kubeconfig file to authenticate to kubernetes api")
}
