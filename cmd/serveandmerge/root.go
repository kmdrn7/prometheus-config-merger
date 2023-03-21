package serveandmerge

import (
	"prometheus-config-merger/pkg/serveandmerge"
	"prometheus-config-merger/pkg/server"
	"prometheus-config-merger/pkg/watcher"

	"github.com/spf13/cobra"
)

var (
	snm *serveandmerge.ServeAndMerge = serveandmerge.NewServeAndMerge()

	dryRun     bool
	kubeconfig string

	listenAddress    string
	reloadUrl        string
	configOutputFile string
)

// Cmd represents the base command when called without any subcommands
var Cmd = &cobra.Command{
	Use:   "serveandmerge",
	Short: "Run Server and Merger in same process.",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		// initialize watcher
		w := watcher.NewWatcher(kubeconfig)

		// initialize server
		s := server.NewServer()
		s.ListenAddress = listenAddress
		s.ReloadUrl = reloadUrl
		s.ConfigOutputFile = configOutputFile

		// create runner list
		runnerList := []serveandmerge.Runner{s, w}

		// run 'em all
		cobra.CheckErr(snm.Start(runnerList))
	},
}

func init() {
	// config var for merger
	Cmd.Flags().BoolVar(&dryRun, "dryrun", false, "Only print watcher configuration information. Usefull for debugging")
	Cmd.Flags().StringVar(&kubeconfig, "kubeconfig", "", "Specify where the kubeconfig file to authenticate to kubernetes api")

	// config var for server
	Cmd.Flags().StringVar(&listenAddress, "web.listen-address", ":5000", "Address on which the webserver service listens")
	Cmd.Flags().StringVar(&reloadUrl, "reload-url", "http://localhost:9090/-/reload", "reload URL to trigger Prometheus reload on")
	Cmd.Flags().StringVar(&configOutputFile, "config-output-file", "/etc/prometheus/prometheus_merged.yaml", "output file from multiple config files merged")
}
