package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/zjyl1994/yusifubot/infra/vars"
)

var app *fiber.App

func Start() {
	app = fiber.New(fiber.Config{DisableStartupMessage: true})
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World 👋!")
	})

	err := app.Listen(vars.ListenAddr)
	if err != nil {
		logrus.Errorln("Http server error", err.Error())
	}
}
