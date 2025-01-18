package http

import "github.com/gofiber/fiber/v2"

func handleIndexPage(c *fiber.Ctx) error {
	return c.SendString("Hello, World 👋!")
}
