package routes

import (
	"todo/src/controller"
	"todo/src/middleware"

	"github.com/gofiber/fiber/v2"
)

func TodoRoutes(app *fiber.App) {
	todo := app.Group("todo", middleware.AuthMiddleware)

	todo.Post("/", controller.CreateTodo)
	todo.Get("/", controller.GetTodos)
	todo.Delete("/:id", controller.DeleteTodo)
	todo.Put("/:id", controller.UpdateTodo)
}
