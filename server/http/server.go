package http

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/zjyl1994/yusifubot/infra/vars"
	"github.com/zjyl1994/yusifubot/service/sqliteadmin"
)

var app *fiber.App

func Start() {
	app = fiber.New(fiber.Config{DisableStartupMessage: true})
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World 👋!")
	})

	sqladmin := sqliteadmin.NewSQLiteAdmin(vars.DBInstance)
	sqladmin.Register(app.Group("/sqladmin"))

	err := app.Listen(vars.ListenAddr)
	if err != nil {
		logrus.Errorln("Http server error", err.Error())
	}
}

func Stop() error {
	if app != nil {
		return app.ShutdownWithTimeout(3 * time.Second)
	}
	return nil
}
