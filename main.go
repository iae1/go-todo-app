package main

import (
	"log"
	"os"

	"go-todo-app/config"
	"go-todo-app/db"
	"go-todo-app/handlers"
	"go-todo-app/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/joho/godotenv"
)

func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = ":3000"
	} else {
		port = ":" + port
	}

	return port
}

func setupRoutes(app *fiber.App, store *session.Store) {
	app.Post("/register", handlers.Register(store))
    app.Post("/login", handlers.Login(store))

	app.Get("/todos", middleware.IsAuthenticated(config.Store), handlers.GetTodos(config.Store))
	app.Post("/todos", middleware.IsAuthenticated(config.Store), handlers.CreateTodo(config.Store))
	app.Put("/todos/:id", middleware.IsAuthenticated(config.Store), handlers.UpdateTodo(config.Store))
	app.Delete("/todos/:id", middleware.IsAuthenticated(config.Store), handlers.DeleteTodo(config.Store))
}

func main() {
	if err := godotenv.Load(); err != nil {
        log.Println("No .env file found")
    }

	config.SetupSessionStore()

	app := fiber.New()

	if err := db.ConnectToDB(); err != nil {
        log.Fatalf("Error connecting to the database: %v", err)
    }
    defer db.CloseDB()

	setupRoutes(app, config.Store)

	log.Fatal(app.Listen(getPort()))
}