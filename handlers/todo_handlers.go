package handlers

import (
	"context"
	"go-todo-app/db"
	"go-todo-app/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

// Helper function to get userId from session
func getUserIdFromSession(c *fiber.Ctx, store *session.Store) (int, error) {
    sess, err := store.Get(c)
    if err != nil {
        return 0, err
    }

    userId := sess.Get("userId")
    if userId == nil {
        return 0, fiber.NewError(fiber.StatusUnauthorized, "user not logged in")
    }

    return userId.(int), nil // Ensure type assertion is safe in your actual implementation
}

func GetTodos(store *session.Store) func(c *fiber.Ctx) error {
    return func(c *fiber.Ctx) error {
        userId, err := getUserIdFromSession(c, store)
        if err != nil {
            return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
        }

        todos := []models.Todo{}
        query := `SELECT id, title, completed FROM todos WHERE user_id = $1`
        rows, err := db.DBPool.Query(context.Background(), query, userId)
        if err != nil {
            return c.Status(500).SendString(err.Error())
        }
        defer rows.Close()

        for rows.Next() {
            var todo models.Todo
            if err := rows.Scan(&todo.ID, &todo.Title, &todo.Completed); err != nil {
                return c.Status(500).SendString(err.Error())
            }
            todos = append(todos, todo)
        }

        return c.JSON(todos)
    }
}

func CreateTodo(store *session.Store) func(c *fiber.Ctx) error {
    return func(c *fiber.Ctx) error {
        userId, err := getUserIdFromSession(c, store)
        if err != nil {
            return err
        }

        todo := new(models.Todo)
        if err := c.BodyParser(todo); err != nil {
            return c.Status(400).SendString(err.Error())
        }

        const insertSQL = `INSERT INTO todos (title, completed, user_id) VALUES ($1, $2, $3) RETURNING id`
        err = db.DBPool.QueryRow(context.Background(), insertSQL, todo.Title, todo.Completed, userId).Scan(&todo.ID)
        if err != nil {
            return c.Status(500).SendString(err.Error())
        }

        return c.Status(201).JSON(todo)
    }
}

func UpdateTodo(store *session.Store) func(c *fiber.Ctx) error {
    return func(c *fiber.Ctx) error {
        userId, err := getUserIdFromSession(c, store)
        if err != nil {
            return err
        }

        id, err := strconv.Atoi(c.Params("id"))
        if err != nil {
            return c.Status(400).SendString(err.Error())
        }

        todo := new(models.Todo)
        if err := c.BodyParser(todo); err != nil {
            return c.Status(400).SendString(err.Error())
        }

        const updateSQL = `UPDATE todos SET title = $1, completed = $2 WHERE id = $3 AND user_id = $4`
        _, err = db.DBPool.Exec(context.Background(), updateSQL, todo.Title, todo.Completed, id, userId)
        if err != nil {
            return c.Status(500).SendString(err.Error())
        }

        todo.ID = id
        return c.JSON(todo)
    }
}

func DeleteTodo(store *session.Store) func(c *fiber.Ctx) error {
    return func(c *fiber.Ctx) error {
        userId, err := getUserIdFromSession(c, store)
        if err != nil {
            return err
        }

        id, err := strconv.Atoi(c.Params("id"))
        if err != nil {
            return c.Status(400).SendString(err.Error())
        }

        const deleteSQL = `DELETE FROM todos WHERE id = $1 AND user_id = $2`
        _, err = db.DBPool.Exec(context.Background(), deleteSQL, id, userId)
        if err != nil {
            return c.Status(500).SendString(err.Error())
        }

        return c.SendStatus(204)
    }
}
