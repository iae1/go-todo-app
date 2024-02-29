package main

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Todo struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

var dbPool *pgxpool.Pool

func main() {
	// Connect to the database
	dbUrl := "postgres://isaaceaston:posgres@localhost:5432/todo_db"
	var err error
	dbPool, err = pgxpool.Connect(context.Background(), dbUrl)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer dbPool.Close()

	app := fiber.New()

	// Setup routes
	app.Get("/todos", getTodos)
	app.Post("/todos", createTodo)
	app.Put("/todos/:id", updateTodo)
	app.Delete("/todos/:id", deleteTodo)

	log.Fatal(app.Listen(":3000"))
}

// Handlers for the routes will go here
