package bot

import (
	"context"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/sirupsen/logrus"
	"github.com/zjyl1994/yusifubot/infra/utils"
	"github.com/zjyl1994/yusifubot/service/catchgame/action"
	"github.com/zjyl1994/yusifubot/service/config"
	"github.com/zjyl1994/yusifubot/service/tg"
)

func Handler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}

	if update.Message.From.IsBot {
		return
	}

	if err := tg.UpdateChatAndUserName(update.Message); err != nil {
		logrus.Warningln("Update chat and user name failed", err.Error())
	}

	err := commandDispatcher(update.Message)
	if err != nil {
		errMsg := "发生错误，请联系管理员"
		if bizErr, ok := err.(utils.BizErr); ok {
			errMsg = bizErr.GetBizMsg()
		} else {
			logrus.Errorln(err)
		}
		utils.ReplyTextToTelegram(update.Message, errMsg, false)
	}
}

func commandDispatcher(msg *models.Message) error {
	command := utils.ParseCommand(msg.Text)
	logrus.Debugln("Received", command, utils.ParseCommandArguments(msg.Text))
	// 在此分发命令
	switch strings.ToLower(command) {
	case "start":
		return utils.ReplyTextToTelegram(msg, "欢迎使用 YusifuBot", false)
	case "debuginfo":
		return tg.InfoHandler(msg)
	case "config":
		return config.Handler(msg)
	case "catch", "catchall", "catch5", "catch10":
		return action.CatchHandler(msg)
	case "catchme":
		return action.CatchMeHandler(msg)
	case "mycatch":
		return action.GetMyCatchHandler(msg)
	case "rankcatch":
		return action.RankCatchHandler(msg)
	case "setmycatch":
		return action.SetMyCatchHandler(msg)
	case "setmyname":
		return action.SetMyNicknameHandler(msg)
	case "setmyemoji":
		return action.SetMyEmojiHandler(msg)
	}
	return nil
}
