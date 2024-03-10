package handlers

import (
	"context"
	"go-todo-app/db"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"golang.org/x/crypto/bcrypt"
)

func Login(store *session.Store) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
        sess, err := store.Get(c)
        if err != nil {
            return c.Status(fiber.StatusInternalServerError).SendString("Cannot get session")
        }

        // Check if user is already logged in
        if sess.Get("username") != nil {
            return c.Status(fiber.StatusBadRequest).SendString("User already logged in")
        }
        
		type LoginRequest struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}

		var request LoginRequest
		if err := c.BodyParser(&request); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
		}

		var passwordHash string
		err := db.DBPool.QueryRow(context.Background(), "SELECT password_hash FROM users WHERE username = $1", request.Username).Scan(&passwordHash)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid login credentials"})
		}

		if err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(request.Password)); err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid login credentials"})
		}

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot get session"})
		}
		sess.Set("username", request.Username)
		if err := sess.Save(); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot save session"})
		}

		return c.SendStatus(fiber.StatusOK)
	}
}

func Register(store *session.Store) func(c *fiber.Ctx) error {
    return func(c *fiber.Ctx) error {
        sess, err := store.Get(c)
        if err != nil {
            return c.Status(fiber.StatusInternalServerError).SendString("Cannot get session")
        }

        // Check if user is already logged in
        if sess.Get("username") != nil {
            return c.Status(fiber.StatusBadRequest).SendString("User already logged in")
        }

        type RegisterRequest struct {
            Username string `json:"username"`
            Password string `json:"password"`
        }
    
        var request RegisterRequest
        if err := c.BodyParser(&request); err != nil {
            return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
        }
    
        // Hash the password
        hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), 14)
        if err != nil {
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot hash password"})
        }
    
        // Insert the user into the database
        _, err = db.DBPool.Exec(context.Background(), "INSERT INTO users (username, password_hash) VALUES ($1, $2)", request.Username, hashedPassword)
        if err != nil {
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot create user"})
        }
    
        return c.SendStatus(fiber.StatusCreated)
    }
}