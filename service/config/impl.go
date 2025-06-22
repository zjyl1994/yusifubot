package config

import (
	"fmt"
	"strings"

	"github.com/go-telegram/bot/models"
	"github.com/zjyl1994/yusifubot/infra/utils"
	"github.com/zjyl1994/yusifubot/infra/vars"
)

func Handler(isGlobal bool) func(msg *models.Message) (err error) {
	return func(msg *models.Message) (err error) {
		var chatId int64
		var msgPrefix string
		if isGlobal {
			chatId = GLOBAL_CONFIG_CHAT_ID
			msgPrefix = "全局"
		} else {
			chatId = msg.Chat.ID
		}

		if msg.From.ID != vars.AdminUserId {
			return utils.ReplyTextToTelegram(msg, "您不是管理员", false)
		}
		commandArgs := utils.ParseCommandArguments(msg.Text)
		switch len(commandArgs) {
		case 0:
			return utils.ReplyTextToTelegram(msg, "请输入配置名", false)
		case 1:
			value, err := Get(chatId, commandArgs[0])
			if err != nil {
				return err
			}
			return utils.ReplyTextToTelegram(msg, fmt.Sprintf("%s配置 '%s' 的值为 '%s'", msgPrefix, commandArgs[0], value), false)
		default:
			data := strings.Join(commandArgs[1:], " ")
			err = Set(chatId, commandArgs[0], data)
			if err != nil {
				return err
			}
			return utils.ReplyTextToTelegram(msg, fmt.Sprintf("%s配置 '%s' 已设置为 '%s'", msgPrefix, commandArgs[0], data), false)
		}
	}
}
