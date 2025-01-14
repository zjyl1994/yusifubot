package catch

import (
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/zjyl1994/yusifubot/infra/utils"
	"github.com/zjyl1994/yusifubot/service/catchgame/common"
	"github.com/zjyl1994/yusifubot/service/catchgame/stamina"
	"github.com/zjyl1994/yusifubot/service/tg"
)

func MyCatch(msg *tgbotapi.Message) error {
	err := tg.UpdateChatAndUserName(msg)
	if err != nil {
		return err
	}
	user := common.UserRel{
		ChatId: msg.Chat.ID,
		UserId: msg.From.ID,
	}

	// 获取体力信息
	sp, err := stamina.GetStaminPoint(user)
	if err != nil {
		return err
	}
	// 获取战利品信息
	loots, err := lootCount(user)
	if err != nil {
		return err
	}
	// 拼装回复信息
	var sb strings.Builder
	sb.WriteString(sp.String())
	sb.WriteString("\n\n*捕捉记录*:\n")
	if len(loots) == 0 {
		sb.WriteString("无\n")
	}
	for _, v := range loots {
		sb.WriteString(fmt.Sprintf("%s x%d\n", utils.EscapeTelegramMarkdown(v.Name), v.Amount))
	}
	return utils.ReplyTextToTelegram(msg, sb.String(), true)
}

// func CatchRank(msg *tgbotapi.Message) error {
// 	err := tg.UpdateChatAndUserName(msg)
// 	if err != nil {
// 		return err
// 	}

// }
