package server

import (
	"prometheus-config-merger/pkg/server"

	"github.com/spf13/cobra"
)

var (
	serverImpl *server.Server = server.NewServer()
)

// Cmd represents the base command when called without any subcommands
var Cmd = &cobra.Command{
	Use:   "server",
	Short: "Merge multiple prometheus config files into single config file.",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		cobra.CheckErr(serverImpl.Start())
	},
}

func init() {
	Cmd.Flags().StringVar(&serverImpl.ListenAddress, "web.listen-address", ":5000", "Address on which the webserver service listens")
	Cmd.Flags().StringVar(&serverImpl.ReloadUrl, "reload-url", "http://localhost:9090/-/reload", "reload URL to trigger Prometheus reload on")
	Cmd.Flags().StringVar(&serverImpl.ConfigOutputFile, "config-output-file", "/etc/prometheus/prometheus_merged.yaml", "output file from multiple config files merged")
}
