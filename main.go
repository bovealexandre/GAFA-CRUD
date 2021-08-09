package main

import (
	"log"
	"os"

	"housecms/App/routes"
	"housecms/config"

	fiber "github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func serveStatic(app *fiber.App) { app.Static("/", "./build") }

func setupRoutes(app *fiber.App) {

	// moved from main method
	api := app.Group("/api")

	routes.TodoRoute(api.Group("/todos"))
}

func main() {
	//connect http
	app := fiber.New()

	//connect logger
	app.Use(logger.New())

	//Handle Cors
	app.Use(cors.New())

	// recover
	app.Use(recover.New())

	config.ConnectDB()

	//Serve the build file
	serveStatic(app)

	//Setup Routes
	setupRoutes(app)

	port := os.Getenv("PORT")

	// Verify if heroku provided the port or not
	if os.Getenv("PORT") == "" {
		port = "5000"
	}

	// Start server on http://${heroku-url}:${port}
	log.Fatal(app.Listen(":" + port))
}
