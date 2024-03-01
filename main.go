package main

import (
	"context"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Todo struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

var dbPool *pgxpool.Pool

func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = ":3000"
	} else {
		port = ":" + port
	}

	return port
}

func main() {
	// Connect to the database
	dbUrl := os.Getenv("DATABASE_URL")
	var err error
	dbPool, err = pgxpool.Connect(context.Background(), dbUrl)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer dbPool.Close()

	app := fiber.New()

	// Setup routess
	app.Get("/todos", getTodos)
	app.Post("/todos", createTodo)
	app.Put("/todos/:id", updateTodo)
	app.Delete("/todos/:id", deleteTodo)

	log.Fatal(app.Listen(getPort()))
}

// Handlers for the routes will go here
