package http

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/zjyl1994/yusifubot/service/catchgame/catch"
	"github.com/zjyl1994/yusifubot/service/catchgame/common"
	"github.com/zjyl1994/yusifubot/service/catchgame/stamina"
	"golang.org/x/sync/singleflight"
)

type catchGameHandler struct {
	statSf singleflight.Group
}

func (h *catchGameHandler) Index(c *fiber.Ctx) error {
	result, err := h.getStatData()
	if err != nil {
		return err
	}

	return c.Render("templates/catchgame_dashboard", fiber.Map{
		"stat": result,
	})
}

func (h *catchGameHandler) getStatData() (map[string][]map[string]any, error) {
	result, err, _ := h.statSf.Do("stat", func() (any, error) {
		return catch.Stat()
	})
	if err != nil {
		return nil, err
	}
	return result.(map[string][]map[string]any), nil
}

func (h *catchGameHandler) Stat(c *fiber.Ctx) error {
	result, err := h.getStatData()
	if err != nil {
		return err
	}
	return c.JSON(result)
}
func (h *catchGameHandler) AddStaminPoint(c *fiber.Ctx) error {
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
