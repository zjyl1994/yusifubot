package tg

import (
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func UpdateChatAndUserName(msg *tgbotapi.Message) error {
	userId := msg.From.ID
	username := GetTgUserName(msg.From)
	if err := UpdateUserName(userId, username); err != nil {
		return err
	}

	chatId := msg.Chat.ID
	chatName := GetTgChatName(msg.Chat)
	return UpdateChatName(chatId, chatName)
}

func GetTgUserName(msg *tgbotapi.User) string {
	return strings.TrimSpace(msg.FirstName + " " + msg.LastName)
}

func GetTgChatName(msg *tgbotapi.Chat) string {
	if msg.Type == "private" {
		return strings.TrimSpace(msg.FirstName + " " + msg.LastName)
	} else {
		return msg.Title
	}
}
