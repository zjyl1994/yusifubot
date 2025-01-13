package utils

import (
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/zjyl1994/yusifubot/infra/vars"
)

const PARSE_MODE_MARKDOWN = "MarkdownV2"

func ReplyTextToTelegram(input *tgbotapi.Message, text string, markdown bool) error {
	msg := tgbotapi.NewMessage(input.Chat.ID, text)
	msg.ReplyToMessageID = input.MessageID
	if markdown {
		msg.ParseMode = PARSE_MODE_MARKDOWN
	}
	_, err := vars.BotInstance.Send(msg)
	return err
}

func ReplyStickerToTelegram(input *tgbotapi.Message, stickerId string) error {
	sticker := tgbotapi.FileID(stickerId)
	msg := tgbotapi.NewSticker(input.Chat.ID, sticker)
	msg.ReplyToMessageID = input.MessageID
	_, err := vars.BotInstance.Send(msg)
	return err
}

func EscapeTelegramMarkdown(input string) string {
	var builder strings.Builder
	for _, char := range input {
		if _, ok := MARKDOWN_ESCAPE_MAP[char]; ok {
			builder.WriteRune('\\')
		}
		builder.WriteRune(char)
	}
	return builder.String()
}

var (
	MARKDOWN_ESCAPE_CHARS = []rune{'_', '*', '[', ']', '(', ')', '~', '`', '>',
		'#', '+', '-', '=', '|', '{', '}', '.', '!'}
	MARKDOWN_ESCAPE_MAP map[rune]struct{}
)

func init() {
	MARKDOWN_ESCAPE_MAP = make(map[rune]struct{})
	for _, r := range MARKDOWN_ESCAPE_CHARS {
		MARKDOWN_ESCAPE_MAP[r] = struct{}{}
	}
}
