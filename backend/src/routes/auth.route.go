package routes

import (
	"todo/src/controller"
	"todo/src/middleware"

	"github.com/gofiber/fiber/v2"
)

func AuthRoutes(app *fiber.App) {
	auth := app.Group("/auth")
	auth.Post("/register", controller.RegisterUser)
	auth.Post("/login", controller.LoginUser)
	auth.Post("/logOut", middleware.AuthMiddleware, controller.LogoutUser)
}
