package config

import (
	"fmt"
	"strings"

	"github.com/go-telegram/bot/models"
	"github.com/zjyl1994/yusifubot/infra/utils"
	"github.com/zjyl1994/yusifubot/infra/vars"
)

func Handler(msg *models.Message) (err error) {
	if msg.From.ID != vars.AdminUserId {
		return utils.ReplyTextToTelegram(msg, "您不是管理员", false)
	}
	commandArgs := utils.ParseCommandArguments(msg.Text)
	switch len(commandArgs) {
	case 0:
		return utils.ReplyTextToTelegram(msg, "请输入配置名", false)
	case 1:
		value, err := Get(msg.Chat.ID, commandArgs[0])
		if err != nil {
			return err
		}
		return utils.ReplyTextToTelegram(msg, fmt.Sprintf("配置 '%s' 的值为 '%s'", commandArgs[0], value), false)
	default:
		data := strings.Join(commandArgs[1:], " ")
		err = Set(msg.Chat.ID, commandArgs[0], data)
		if err != nil {
			return err
		}
		return utils.ReplyTextToTelegram(msg, fmt.Sprintf("配置 '%s' 已设置为 '%s'", commandArgs[0], data), false)
	}
}
