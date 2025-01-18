package http

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/zjyl1994/yusifubot/service/catchgame/common"
	"github.com/zjyl1994/yusifubot/service/catchgame/stamina"
)

type catchGameHandler struct{}

func (catchGameHandler) Index(c *fiber.Ctx) error {
	return c.SendString("index page")
}

func (catchGameHandler) AddStaminPoint(c *fiber.Ctx) error {
	userId, err := strconv.ParseInt(c.FormValue("user"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("user:" + err.Error())
	}
	chatId, err := strconv.ParseInt(c.FormValue("chat"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("chat:" + err.Error())
	}
	amount, err := strconv.ParseInt(c.FormValue("amount"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("amount:" + err.Error())
	}
	if userId == 0 || chatId == 0 || amount == 0 {
		return c.Status(fiber.StatusBadRequest).SendString("user, chat, amount must be set")
	}
	userRel := common.UserRel{ChatId: chatId, UserId: userId}
	err = stamina.AddStaminPoint(userRel, amount)
	if err != nil {
		return err
	}
	sp, err := stamina.GetStaminPoint(userRel)
	if err != nil {
		return err
	}
	return c.SendString(sp.String())
}
