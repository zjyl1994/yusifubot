package utils

import (
	"context"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/vinta/pangu"
	"github.com/zjyl1994/yusifubot/infra/vars"
)

func ReplyTextToTelegram(input *models.Message, text string, markdown bool) error {
	var msgParams bot.SendMessageParams
	msgParams.Text = text
	msgParams.ChatID = input.Chat.ID
	msgParams.DisableNotification = true
	msgParams.ReplyParameters = &models.ReplyParameters{
		MessageID: input.ID,
	}
	if markdown {
		msgParams.ParseMode = models.ParseModeMarkdown
	} else {
		msgParams.Text = pangu.SpacingText(msgParams.Text)
	}
	_, err := vars.BotInstance.SendMessage(context.Background(), &msgParams)
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

func ParseCommand(text string) string {
	// 如果输入为空则返回空
	if text == "" {
		return ""
	}

	// 分割文本为单词数组
	parts := strings.Fields(text)
	if len(parts) == 0 {
		return ""
	}

	// 获取命令（去除可能的@机器人名称）
	command := parts[0]
	if strings.Contains(command, "@") {
		command = strings.Split(command, "@")[0]
	}
	// 去除命令前的/
	command = strings.TrimPrefix(command, "/")

	return command
}

func ParseCommandArguments(text string) []string {
	// 如果输入为空则返回空切片
	if text == "" {
		return []string{}
	}

	// 分割文本为单词数组
	parts := strings.Fields(text)
	if len(parts) == 0 {
		return []string{}
	}

	// 如果只有命令没有参数，返回空切片
	if len(parts) == 1 {
		return []string{}
	}

	// 返回命令后的所有参数
	return parts[1:]
}

func IsGroup(msg *models.Message) bool {
	return msg.Chat.Type == models.ChatTypeGroup || msg.Chat.Type == models.ChatTypeSupergroup
}
