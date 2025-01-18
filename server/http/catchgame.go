package http

import "github.com/gofiber/fiber/v2"

type catchGameHandler struct{}

func (catchGameHandler) Index(c *fiber.Ctx) error {
	return c.SendString("index page")
}
