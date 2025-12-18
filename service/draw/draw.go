package draw

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/replicate/replicate-go"
	"github.com/sirupsen/logrus"
	"github.com/zjyl1994/yusifubot/infra/utils"
	"github.com/zjyl1994/yusifubot/infra/vars"
)

const (
	MODEL_IDENTIFIER = "prunaai/z-image-turbo:7ea16386290ff5977c7812e66e462d7ec3954d8e007a8cd18ded3e7d41f5d7cf"
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
		url, err := drawWithReplicate(ctx, prompt)
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
	output, err := r8.Run(ctx, MODEL_IDENTIFIER, replicate.PredictionInput{
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
