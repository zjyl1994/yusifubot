package bot

import (
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
	"github.com/zjyl1994/yusifubot/infra/utils"
	"github.com/zjyl1994/yusifubot/infra/vars"
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

	// TODO: 在此分发命令

	return utils.ReplyTextToTelegram(msg, "未知命令", false)
}
