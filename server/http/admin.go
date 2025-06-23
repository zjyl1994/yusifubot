package http

import (
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/zjyl1994/yusifubot/infra/vars"
	"github.com/zjyl1994/yusifubot/service/catchgame/catchobj"
	"github.com/zjyl1994/yusifubot/service/catchgame/catchret"
	"github.com/zjyl1994/yusifubot/service/catchgame/common"
	"github.com/zjyl1994/yusifubot/service/catchgame/stamina"
	"gorm.io/gorm"
)

func AdminApi(app *fiber.App) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("YUSIFUBOT")
	})
	app.Get("/catchdata", getCatchData)
	adminGroup := app.Group("/admin", auth)
	adminGroup.Post("/givesp", giveSp)
	adminGroup.Post("/givecatch", giveCatch)
	adminGroup.Post("/createobj", createCatchObj)
}

func auth(c *fiber.Ctx) error {
	authHeader := c.Get(fiber.HeaderAuthorization)
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).SendString("missing authorization header")
	}

	// 分割 "Bearer" 和 token
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return c.Status(fiber.StatusUnauthorized).SendString("invalid authorization format")
	}

	// 验证 token
	if parts[1] != vars.AdminToken {
		return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
	}
	return c.Next()
}

func giveSp(c *fiber.Ctx) error {
	var req []struct {
		ChatId int64 `json:"chat_id"`
		UserId int64 `json:"user_id"`
		Num    int64 `json:"num"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("invalid request body:" + err.Error())
	}

	err := vars.DBInstance.Transaction(func(tx *gorm.DB) error {
		for _, item := range req {
			if item.Num <= 0 {
				return c.Status(fiber.StatusBadRequest).SendString("num must be greater than 0")
			}
			if err := stamina.AddStaminPoint(tx, common.UserRel{
				ChatId: item.ChatId,
				UserId: item.UserId,
			}, item.Num); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	return c.SendString("all success")
}

func giveCatch(c *fiber.Ctx) error {
	var req []struct {
		ChatId int64 `json:"chat_id"`
		UserId int64 `json:"user_id"`
		ObjId  int64 `json:"obj_id"`
		Num    int64 `json:"num"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("invalid request body:" + err.Error())
	}

	err := vars.DBInstance.Transaction(func(tx *gorm.DB) error {
		for _, item := range req {
			if err := catchret.GiveCatchNum(tx, common.UserRel{
				ChatId: item.ChatId,
				UserId: item.UserId,
			}, item.ObjId, item.Num); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	return c.SendString("all success")
}

func createCatchObj(c *fiber.Ctx) error {
	var req []struct {
		ChatId int64  `json:"chat_id"`
		UserId int64  `json:"user_id"`
		Name   string `json:"name"`
		Emoji  string `json:"emoji"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("invalid request body:" + err.Error())
	}

	err := vars.DBInstance.Transaction(func(tx *gorm.DB) error {
		for _, item := range req {
			if err := catchobj.CreateCatchObj(tx, common.UserRel{
				ChatId: item.ChatId,
				UserId: item.UserId,
			}, item.Name, item.Emoji); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	return c.SendString("all success")
}

func getCatchData(c *fiber.Ctx) error {
	var chatId, userId int64

	if str := c.Query("chat_id"); len(str) > 0 {
		if i64, err := strconv.ParseInt(str, 10, 64); err == nil {
			chatId = i64
		} else {
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}
	} else {
		return c.Status(fiber.StatusBadRequest).SendString("chat_id is required")
	}
	if str := c.Query("user_id"); len(str) > 0 {
		if i64, err := strconv.ParseInt(str, 10, 64); err == nil {
			userId = i64
		} else {
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}
	} else {
		return c.Status(fiber.StatusBadRequest).SendString("user_id is required")
	}
	result, err := catchret.GetMyCatch(vars.DBInstance, common.UserRel{ChatId: chatId, UserId: userId})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	return c.JSON(result)
}
