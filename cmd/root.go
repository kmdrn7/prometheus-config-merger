package cmd

import (
	"log"
	"os"

	"prometheus-config-merger/cmd/merge"
	"prometheus-config-merger/cmd/server"
	"prometheus-config-merger/pkg/config"
	"prometheus-config-merger/pkg/http"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	debug   bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "prometheus-config-merger",
	Short: "Merge multiple prometheus config files into single config file.",
	Long: `
Merge rule_files and scrape_configs across multiple prometheus config files into single config file.
Just in case you want to segregate prometheus configuration into multiple file, and have it merged before consumed by prometheus-server.
`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	cobra.OnInitialize(http.New)

	// add merge subcommand
	rootCmd.AddCommand(merge.MergeCmd)

	// add server subcommand
	rootCmd.AddCommand(server.ServerCmd)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "Config file (default is ./prometheus-config-merger.yaml)")
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "Enable debug log message")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Search config in current directory with name "prometheus-config-merger.yaml".
		viper.AddConfigPath(".")
		viper.SetConfigType("yaml")
		viper.SetConfigName("prometheus-config-merger")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// override configuration through viper
	viper.Set("debug", debug)
	if debug {
		log.Println("running in DEBUG mode")
	}

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.Println("using config file:", viper.ConfigFileUsed())

		pcmConf := &config.Config{}
		if err := viper.Unmarshal(pcmConf); err != nil {
			log.Fatal(err.Error())
		}
	}
}
