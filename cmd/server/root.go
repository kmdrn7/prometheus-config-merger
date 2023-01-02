package server

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/cobra"
)

var ServerCmd = &cobra.Command{
	Use:   "server",
	Short: "Merge multiple prometheus config files into single config file.",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		app := fiber.New(fiber.Config{
			DisableStartupMessage: true,
		})

		app.Get("/healthz", func(c *fiber.Ctx) error {
			return c.SendStatus(200)
		})

		log.Fatal(app.Listen(":5000"))
	},
}
