package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/zjyl1994/yusifubot/infra/vars"
)

var app *fiber.App

func Start() {
	app = fiber.New(fiber.Config{DisableStartupMessage: true})
	app.Get("/", handleIndexPage)
	app.Get("/ntunnel", ntunnelPage)
	app.Post("/ntunnel", ntunnelGen())

	err := app.Listen(vars.ListenAddr)
	if err != nil {
		logrus.Errorln("Http server error", err.Error())
	}
}
