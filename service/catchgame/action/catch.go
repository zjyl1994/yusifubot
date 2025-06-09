package action

import (
	"errors"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func CatchHandler(msg *tgbotapi.Message) error {
	if msg.ReplyToMessage != nil { // 捕捉特定用户
		return catchSpecificHandler(msg)
	} else { // 捕捉混池
		return catchPoolHandler(msg)
	}
}

// 只捉特定用户
func catchSpecificHandler(msg *tgbotapi.Message) error {
	return errors.New("TODO")
}

// 混池捕捉
func catchPoolHandler(msg *tgbotapi.Message) error {
	return errors.New("TODO")
}
