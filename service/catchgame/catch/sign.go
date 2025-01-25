package catch

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/zjyl1994/yusifubot/infra/utils"
	"github.com/zjyl1994/yusifubot/service/catchgame/common"
	"github.com/zjyl1994/yusifubot/service/catchgame/sign"
	"github.com/zjyl1994/yusifubot/service/tg"
)

func SignAction(msg *tgbotapi.Message) error {
	err := tg.UpdateChatAndUserName(msg)
	if err != nil {
		return err
	}

	user := common.UserRel{
		ChatId: msg.Chat.ID,
		UserId: msg.From.ID,
	}

	awardSp, err := sign.Sign(user)
	if err != nil {
		return err
	}

	return utils.ReplyTextToTelegram(msg, fmt.Sprintf("签到成功,奖励体力%d点", awardSp), false)
}
