package src

import (
	"log"
	"todo/src/db"
	"todo/src/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func SetupApp() *fiber.App {
	app := fiber.New()
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loging to .env file")
	}
	db.ConnectDB()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Welcome to fiber app")
	})
	routes.AuthRoutes(app)
	return app
}
