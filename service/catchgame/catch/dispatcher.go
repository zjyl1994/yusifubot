package catch

import (
	"regexp"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/zjyl1994/yusifubot/infra/utils"
	"github.com/zjyl1994/yusifubot/service/tg"
)

var catchCommandRegexp = regexp.MustCompile(`(?i)^catch([a-zA-Z]+?)(\d+|all)?$`)

// 所有catch开头命令由此分发
func CatchDispatcher(msg *tgbotapi.Message) error {
	command := msg.Command()
	args := strings.Fields(msg.CommandArguments())

	err := tg.UpdateChatAndUserName(msg)
	if err != nil {
		return err
	}
	// 单纯的捕捉指令，结构化响应
	if command == "catch" {
		switch len(args) {
		case 0:
			return utils.ReplyTextToTelegram(msg, "需要指定捕捉对象", false)
		case 1:
			return CatchAction(msg, args[0], catchNum("1"))
		default:
			return CatchAction(msg, args[0], catchNum(args[1]))
		}
	}
	// 组合指令，尝试使用正则解析
	if matches := catchCommandRegexp.FindStringSubmatch(command); matches != nil {
		var num catchNum
		catchName := matches[1]

		if len(matches) > 2 && matches[2] != "" {
			num = catchNum(matches[2])
		} else {
			num = catchNum("1")
		}

		return CatchAction(msg, catchName, num)
	}
	return utils.ReplyTextToTelegram(msg, "无法解析命令"+command, false)
}
