package http

import (
	"database/sql"
	"embed"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/template/html/v2"
	"github.com/sirupsen/logrus"
	"github.com/zjyl1994/yusifubot/infra/vars"
	"github.com/zjyl1994/yusifubot/service/ntunnel"

	sqlite3 "modernc.org/sqlite/lib"
)

var app *fiber.App

//go:embed templates/*
var tmpl embed.FS

func Start() {
	catchGameHandler := new(catchGameHandler)

	engine := html.NewFileSystem(http.FS(tmpl), ".html")

	app = fiber.New(fiber.Config{
		DisableStartupMessage: true,
		Views:                 engine,
	})
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("It works!")
	})

	if vars.AdminUser != "" && vars.AdminPass != "" {
		// 远程DB隧道
		tun := ntunnel.NewNTunnel(
			func(s string) (*sql.DB, error) {
				return vars.DBInstance.DB()
			},
			func(err error) (int, string) {
				if err != nil {
					return 1, err.Error()
				}
				return 0, ""
			},
			sqlite3.SQLITE_VERSION)
		// 登录组件
		authCfg := basicauth.ConfigDefault
		authCfg.Users = map[string]string{
			vars.AdminUser: vars.AdminPass,
		}
		authWare := basicauth.New(authCfg)
		// 登录限流
		limiterCfg := limiter.ConfigDefault
		limiterCfg.Max = 30
		limiterCfg.KeyGenerator = func(c *fiber.Ctx) string {
			if val := c.Get("X-Real-Ip"); val != "" {
				return val
			}
			return c.IP()
		}
		rateLimiter := limiter.New(limiterCfg)
		// 路由组
		adminG := app.Group("/admin", rateLimiter, authWare)
		adminG.Post("/addsp", catchGameHandler.AddStaminPoint) // 增加体力API
		adminG.Post("/ntunnel", adaptor.HTTPHandler(tun))      // 远程DB隧道
	}

	catchG := app.Group("/catchgame")
	catchG.Get("/", catchGameHandler.Index)
	catchG.Get("/stat", catchGameHandler.Stat)

	err := app.Listen(vars.ListenAddr)
	if err != nil {
		logrus.Errorln("Http server error", err.Error())
	}
}
