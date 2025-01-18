package http

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/zjyl1994/yusifubot/infra/vars"
	"github.com/zjyl1994/yusifubot/service/ntunnel"
	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
)

func ntunnelGen() fiber.Handler {
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
	return adaptor.HTTPHandler(tun)
}

func ntunnelPage(c *fiber.Ctx) error {
	return c.SendString("Hello, World 👋!")
}
