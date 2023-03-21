package server

import (
	"log"
	"prometheus-config-merger/pkg/config"
	"prometheus-config-merger/pkg/http"
	"prometheus-config-merger/pkg/merge"
	"prometheus-config-merger/pkg/utils"
	"prometheus-config-merger/pkg/watcher"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

type Server struct {
	ReloadUrl        string
	ListenAddress    string
	ConfigOutputFile string
}

func NewServer() *Server {
	return &Server{
		ReloadUrl:        "",
		ListenAddress:    "",
		ConfigOutputFile: "",
	}
}

func (s *Server) Start() error {

	cfg := config.GetConfig()
	log.Println("using prometheus server at", cfg.PrometheusEndpoint)

	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})

	// add middleware to log requests
	app.Use(logger.New(logger.ConfigDefault))

	// add middleware to recover from panic
	app.Use(recover.New())

	// implement readiness/healthcheck
	app.Get("/healthz", func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})

	// forward reload request from sidecar to prometheus-server
	app.Post("/-/reload", func(c *fiber.Ctx) error {
		log.Println("merging prometheus config files")
		merge.RunWithReturn(cfg)

		log.Println("trigger reload on prometheus-server")
		if r, err := http.Post(s.ReloadUrl); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		} else {
			log.Println(r.StatusCode())
		}

		return c.SendStatus(200)
	})

	app.Post("/-/inject-config", func(c *fiber.Ctx) error {
		payload := &watcher.ResourceHTTPBodyPayload{}
		if err := c.BodyParser(payload); err != nil {
			panic(err)
		}
		log.Printf("[SERVER] syncing content to file %s \n", payload.Filepath)
		if err := utils.SyncResourceContentToLocalFile([]byte(payload.Content), payload.Filepath); err != nil {
			panic(err)
		}

		return c.SendStatus(200)
	})

	log.Println("starting webserver on", s.ListenAddress)
	return app.Listen(s.ListenAddress)
}
