package main

import (
	"context"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func getTodos(c *fiber.Ctx) error {
	todos := []Todo{}
	rows, err := dbPool.Query(context.Background(), "SELECT id, title, completed FROM todoz")
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		var todo Todo
		if err := rows.Scan(&todo.ID, &todo.Title, &todo.Completed); err != nil {
			return c.Status(500).SendString(err.Error())
		}
		todos = append(todos, todo)
	}

	return c.JSON(todos)
}

func createTodo(c *fiber.Ctx) error {
	todo := new(Todo)

	if err := c.BodyParser(todo); err != nil {
		return c.Status(400).SendString(err.Error())
	}

	const insertSQL = `INSERT INTO todoz (title, completed) VALUES ($1, $2) RETURNING id`
	err := dbPool.QueryRow(context.Background(), insertSQL, todo.Title, todo.Completed).Scan(&todo.ID)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	return c.Status(201).JSON(todo)
}

func updateTodo(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).SendString(err.Error())
	}

	todo := new(Todo)

	if err := c.BodyParser(todo); err != nil {
		return c.Status(400).SendString(err.Error())
	}

	const updateSQL = `UPDATE todoz SET title = $1, completed = $2 WHERE id = $3`
	_, err = dbPool.Exec(context.Background(), updateSQL, todo.Title, todo.Completed, id)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	todo.ID = id
	return c.JSON(todo)
}

func deleteTodo(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).SendString(err.Error())
	}

	const deleteSQL = `DELETE FROM todoz WHERE id = $1`
	_, err = dbPool.Exec(context.Background(), deleteSQL, id)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	return c.SendStatus(204)
}
