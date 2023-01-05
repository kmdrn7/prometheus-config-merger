package server

import (
	"log"
	"prometheus-config-merger/pkg/config"
	"prometheus-config-merger/pkg/http"
	"prometheus-config-merger/pkg/merge"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/spf13/cobra"
)

var (
	listenAddress    string
	reloadUrl        string
	configOutputFile string
)

var ServerCmd = &cobra.Command{
	Use:   "server",
	Short: "Merge multiple prometheus config files into single config file.",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		cfg := config.GetConfig()
		log.Println("using prometheus server at", cfg.PrometheusEndpoint)

		app := fiber.New(fiber.Config{
			DisableStartupMessage: true,
		})

		// add middleware to log requests
		app.Use(logger.New(logger.ConfigDefault))

		// implement readiness/healthcheck
		app.Get("/healthz", func(c *fiber.Ctx) error {
			return c.SendStatus(200)
		})

		// forward reload request from sidecar to prometheus-server
		app.Post("/-/reload", func(c *fiber.Ctx) error {
			log.Println("merging prometheus config files")
			merge.Run()

			log.Println("trigger reload on prometheus-server")
			if r, err := http.Post(reloadUrl); err != nil {
				log.Panic(err.Error())
			} else {
				log.Println(r.StatusCode())
			}

			return c.SendStatus(200)
		})

		log.Println("starting webserver on", listenAddress)
		log.Fatal(app.Listen(listenAddress))
	},
}

func init() {
	ServerCmd.Flags().StringVar(&listenAddress, "web.listen-address", ":5000", "Address on which the webserver service listens")
	ServerCmd.Flags().StringVar(&reloadUrl, "reload-url", "http://localhost:9090/-/reload", "reload URL to trigger Prometheus reload on")
	ServerCmd.Flags().StringVar(&configOutputFile, "config-output-file", "/etc/prometheus/prometheus_merged.yaml", "output file from multiple config files merged")
}
