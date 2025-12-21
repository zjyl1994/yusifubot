package draw

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/replicate/replicate-go"
	"github.com/sirupsen/logrus"
	"github.com/zjyl1994/yusifubot/infra/utils"
	"github.com/zjyl1994/yusifubot/infra/vars"
)

const (
	IMAGE_MODEL_IDENTIFIER  = "prunaai/z-image-turbo"
	PROMPT_MODEL_IDENTIFIER = "openai/gpt-5-nano"
)

func DrawImageHandler(msg *models.Message) error {
	if msg.From.ID != vars.AdminUserId && !checkChatIdAllowed(msg.Chat.ID) {
		return utils.ReplyTextToTelegram(msg, "该聊天未开启绘图能力", false)
	}
	commandArgs := utils.ParseCommandArguments(msg.Text)
	prompt := strings.TrimSpace(strings.Join(commandArgs, " "))
	if prompt == "" {
		return utils.ReplyTextToTelegram(msg, "请输入提示词", false)
	}

	if msg.From.ID != vars.AdminUserId {
		key := "draw:" + strconv.FormatInt(msg.From.ID, 10)
		if !vars.ReplicateCooldown.CheckAndSetCooldown(key, 1*time.Minute) {
			remaining := vars.ReplicateCooldown.RemainingTime(key)
			return utils.ReplyTextToTelegram(msg, fmt.Sprintf("你的绘图能力冷却中，请等待 %d 秒", int(remaining.Seconds())), false)
		}
	}

	go func(ctx context.Context, m *models.Message) {
		preparedPrompt, err := preparePrompt(context.Background(), prompt)
		if err != nil {
			utils.ReplyTextToTelegram(msg, err.Error(), false)
			return
		}
		url, err := drawWithReplicate(ctx, preparedPrompt)
		if err != nil {
			utils.ReplyTextToTelegram(msg, err.Error(), false)
			return
		}

		params := &bot.SendPhotoParams{
			ChatID: msg.Chat.ID,
			ReplyParameters: &models.ReplyParameters{
				MessageID: msg.ID,
			},
			Photo:               &models.InputFileString{Data: url},
			DisableNotification: true,
		}

		_, err = vars.BotInstance.SendPhoto(ctx, params)
		if err != nil {
			utils.ReplyTextToTelegram(msg, err.Error(), false)
			return
		}
	}(context.Background(), msg)
	return nil
}

func drawWithReplicate(ctx context.Context, prompt string) (string, error) {
	r8, err := replicate.NewClient(replicate.WithToken(vars.ReplicateToken))
	if err != nil {
		return "", utils.NewBizErr("create replicate client failed")
	}
	output, err := r8.Run(ctx, IMAGE_MODEL_IDENTIFIER, replicate.PredictionInput{
		"prompt": prompt,
	}, nil)
	if err != nil {
		return "", utils.NewBizErr("run replicate model failed")
	}
	logrus.Debugf("replicate output: %s", utils.MarshalToJsonNoError(output))
	outputURL, ok := output.(string)
	if !ok {
		return "", utils.NewBizErr("output is not a string")
	}
	return outputURL, nil
}

func preparePrompt(ctx context.Context, prompt string) (string, error) {
	prompt = strings.TrimSpace(prompt)
	if prompt == "" {
		return "", utils.NewBizErr("prompt is empty")
	}
	if utf8.RuneCountInString(prompt) > 50 { // 本身就很长的prompt不做拓展
		return prompt, nil
	}
	// 使用gpt-5-nano进行拓展
	r8, err := replicate.NewClient(replicate.WithToken(vars.ReplicateToken))
	if err != nil {
		return "", utils.NewBizErr("create replicate client failed")
	}
	output, err := r8.Run(ctx, PROMPT_MODEL_IDENTIFIER, replicate.PredictionInput{
		"system_prompt":         "你是一个画师，解读输入的文字，转化为尽可能明确的绘图指令并输出，无需输出更多其他内容。",
		"prompt":                prompt,
		"max_completion_tokens": 1024,
	}, nil)
	if err != nil {
		return "", utils.NewBizErr("run replicate model failed")
	}
	logrus.Debugf("replicate output: %s", utils.MarshalToJsonNoError(output))

	var outputText string
	if str, ok := output.(string); ok {
		outputText = str
	} else if tokens, ok := output.([]interface{}); ok {
		sb := strings.Builder{}
		for _, token := range tokens {
			if s, ok := token.(string); ok {
				sb.WriteString(s)
			}
		}
		outputText = sb.String()
	} else {
		return "", utils.NewBizErr("output format error")
	}
	return outputText, nil
}
