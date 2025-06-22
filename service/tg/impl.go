package tg

import (
	"fmt"
	"strings"

	"github.com/go-telegram/bot/models"
	"github.com/zjyl1994/yusifubot/infra/utils"
)

func UpdateChatAndUserName(msg *models.Message) error {
	userId := msg.From.ID
	username := GetTgUserName(msg.From)
	if err := UpdateUserName(userId, username); err != nil {
		return err
	}

	chatId := msg.Chat.ID
	chatName := GetTgChatName(&msg.Chat)
	return UpdateChatName(chatId, chatName)
}

func GetTgUserName(msg *models.User) string {
	return strings.TrimSpace(msg.FirstName + " " + msg.LastName)
}

func GetTgChatName(msg *models.Chat) string {
	if msg.Type == "private" {
		return strings.TrimSpace(msg.FirstName + " " + msg.LastName)
	} else {
		return msg.Title
	}
}

func InfoHandler(msg *models.Message) error {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("ChatId: %d\n", msg.Chat.ID))
	sb.WriteString(fmt.Sprintf("UserId: %d\n", msg.From.ID))
	if msg.ReplyToMessage != nil {
		sb.WriteString(fmt.Sprintf("ReplyUserId: %d\n", msg.ReplyToMessage.From.ID))
	}
	return utils.ReplyTextToTelegram(msg, sb.String(), false)
}
