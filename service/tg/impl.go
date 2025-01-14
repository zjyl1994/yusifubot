package tg

import (
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func UpdateChatAndUserName(msg *tgbotapi.Message) error {
	userId := msg.From.ID
	username := strings.TrimSpace(msg.From.FirstName + " " + msg.From.LastName)
	if err := UpdateUserName(userId, username); err != nil {
		return err
	}

	chatId := msg.Chat.ID
	var chatName string
	if msg.Chat.Type == "private" {
		chatName = strings.TrimSpace(msg.Chat.FirstName + " " + msg.Chat.LastName)
	} else {
		chatName = msg.Chat.Title
	}
	return UpdateChatName(chatId, chatName)
}
