package tg

import (
	"strings"

	"github.com/go-telegram/bot/models"
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
