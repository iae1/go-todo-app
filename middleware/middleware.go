package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

func IsAuthenticated(store *session.Store) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		sess, err := store.Get(c)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
		}

		if sess.Get("userID") == nil {
			return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
		}

		return c.Next()
	}
}
