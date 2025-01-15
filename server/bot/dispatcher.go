package bot

import (
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
	"github.com/zjyl1994/yusifubot/infra/utils"
	"github.com/zjyl1994/yusifubot/infra/vars"
	"github.com/zjyl1994/yusifubot/service/catchgame/catch"
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

	// catch开头的命令逻辑复杂需要单独分发逻辑处理
	if strings.HasPrefix(command, "catch") {
		return catch.CatchDispatcher(msg)
	}
	// 在此分发其他命令
	switch strings.ToLower(command) {
	case "debug":
		return handleDebugInfo(msg)
	case "mycatch":
		return catch.MyCatch(msg)
	case "rankcatch":
		return catch.CatchRank(msg)
	default:
		return utils.ReplyTextToTelegram(msg, "未知命令", false)
	}
}
