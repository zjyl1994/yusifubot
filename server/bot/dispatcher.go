package bot

import (
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
	"github.com/zjyl1994/yusifubot/infra/utils"
	"github.com/zjyl1994/yusifubot/infra/vars"
	"github.com/zjyl1994/yusifubot/service/catchgame/action"
	"github.com/zjyl1994/yusifubot/service/tg"
)

func Start() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := vars.BotInstance.GetUpdatesChan(u)

	logrus.Infoln("Bot started")

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if !update.Message.IsCommand() {
			continue
		}

		if update.Message.From.IsBot {
			continue
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
}

func commandDispatcher(msg *tgbotapi.Message) error {
	command := msg.Command()
	args := strings.Fields(msg.CommandArguments())
	logrus.Debugln("Received", command, args)
	if err := tg.UpdateChatAndUserName(msg); err != nil {
		logrus.Warningln("Update chat and user name failed", err.Error())
	}
	// 在此分发命令
	switch strings.ToLower(command) {
	case "start":
		return utils.ReplyTextToTelegram(msg, "欢迎使用 YusifuBot", false)
	case "catch", "catchall", "catch5", "catch10":
		return action.CatchHandler(msg)
	case "catchme":
		return action.CatchMeHandler(msg)
	case "mycatch":

	case "rankcatch":

	case "setnickname":
		return action.SetMyNicknameHandler(msg)
	case "setemoji":
		return action.SetMyEmojiHandler(msg)
	}
	return nil
}
