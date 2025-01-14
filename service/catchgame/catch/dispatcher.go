package catch

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/zjyl1994/yusifubot/infra/utils"
	"github.com/zjyl1994/yusifubot/service/tg"
)

var catchCommandRegexp = regexp.MustCompile(`(?i)^catch([a-zA-Z]+?)(\d+|all)?$`)

// 所有catch开头命令由此分发
func CatchDispatcher(msg *tgbotapi.Message) error {
	command := msg.Command()
	args := msg.CommandArguments()

	err := tg.UpdateChatAndUserName(msg)
	if err != nil {
		return err
	}
	// 单纯的捕捉指令，合并成组合指令走正则解析
	if command == "catch" {
		if len(args) == 0 {
			return utils.ReplyTextToTelegram(msg, "👀 你要捉谁？", false)
		}
		var builder strings.Builder
		builder.WriteString(command)
		for _, r := range args {
			if !unicode.IsSpace(r) {
				builder.WriteRune(r)
			}
		}
		command = builder.String()
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
	return utils.ReplyTextToTelegram(msg, fmt.Sprintf("无法解析：%s %s", command, args), false)
}
