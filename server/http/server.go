package http

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/sirupsen/logrus"
	"github.com/zjyl1994/yusifubot/infra/vars"
	"github.com/zjyl1994/yusifubot/service/ntunnel"
	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
)

var app *fiber.App

func Start() {
	app = fiber.New(fiber.Config{DisableStartupMessage: true})
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Redirect("/catchgame")
	})

	if vars.AdminUser != "" && vars.AdminPass != "" {
		tun := ntunnel.NewNTunnel(
			func(s string) (*sql.DB, error) {
				return vars.DBInstance.DB()
			},
			func(err error) (int, string) {
				if err != nil {
					if sqliteErr, ok := err.(*sqlite.Error); ok {
						return sqliteErr.Code(), sqliteErr.Error()
					} else {
						return 1, err.Error()
					}
				}
				return 0, ""
			},
			sqlite3.SQLITE_VERSION)

		authCfg := basicauth.ConfigDefault
		authCfg.Users = map[string]string{
			vars.AdminUser: vars.AdminPass,
		}
		authWare := basicauth.New(authCfg)

		limiterCfg := limiter.ConfigDefault
		limiterCfg.KeyGenerator = func(c *fiber.Ctx) string {
			if val := c.Get("X-Real-Ip"); val != "" {
				return val
			}
			return c.IP()
		}
		rateLimiter := limiter.New(limiterCfg)

		adminG := app.Group("/admin", rateLimiter, authWare)
		adminG.Post("/ntunnel", adaptor.HTTPHandler(tun))
	}

	var catchGameHandler catchGameHandler
	catchG := app.Group("/catchgame")
	catchG.Get("/", catchGameHandler.Index)

	err := app.Listen(vars.ListenAddr)
	if err != nil {
		logrus.Errorln("Http server error", err.Error())
	}
}
