package bot

import (
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/zjyl1994/yusifubot/infra/utils"
	"github.com/zjyl1994/yusifubot/service/tg"
)

func handleDebugInfo(msg *tgbotapi.Message) error {
	err := tg.UpdateChatAndUserName(msg)
	if err != nil {
		return err
	}
	var sb strings.Builder
	sb.WriteString("Debug Info:\n")

	sb.WriteString("ChatId: ")
	sb.WriteString(strconv.FormatInt(msg.Chat.ID, 10))
	sb.WriteRune('\n')

	sb.WriteString("Chat:")
	sb.WriteString(tg.GetTgChatName(msg.Chat))
	sb.WriteRune('\n')

	sb.WriteString("UserId: ")
	sb.WriteString(strconv.FormatInt(msg.From.ID, 10))
	sb.WriteRune('\n')

	sb.WriteString("User:")
	sb.WriteString(tg.GetTgUserName(msg.From))
	sb.WriteRune('\n')

	return utils.ReplyTextToTelegram(msg, sb.String(), false)
}
