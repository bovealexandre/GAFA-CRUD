package routes

import (
	"housecms/App/controllers"

	"github.com/gofiber/fiber/v2"
)

func TodoRoute(route fiber.Router) {
	route.Get("", controllers.GetTodos)
	route.Post("", controllers.CreateTodo)
	route.Put("/:id", controllers.UpdateTodo)    // new
	route.Delete("/:id", controllers.DeleteTodo) // new
	route.Get("/:id", controllers.GetTodo)       // new
}
