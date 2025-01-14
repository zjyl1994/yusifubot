package catch

import (
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/zjyl1994/yusifubot/infra/utils"
	"github.com/zjyl1994/yusifubot/service/catchgame/catchobj"
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
	loots, err := userLootCount(user)
	if err != nil {
		return err
	}
	// 拼装回复信息
	var sb strings.Builder
	sb.WriteString(sp.String())
	sb.WriteString("\n\n*捕捉记录*:\n")
	if len(loots) == 0 {
		sb.WriteString("无\n")
	} else {
		for _, v := range loots {
			sb.WriteString(fmt.Sprintf("%s x%d\n", utils.EscapeTelegramMarkdown(v.Name), v.Amount))
		}
	}
	return utils.ReplyTextToTelegram(msg, sb.String(), true)
}

func CatchRank(msg *tgbotapi.Message) error {
	err := tg.UpdateChatAndUserName(msg)
	if err != nil {
		return err
	}

	var objId int64
	var objName string
	var objEmoji string
	if target := msg.CommandArguments(); len(target) > 0 {
		cobj, err := catchobj.GetCatchObjByShorthand(target)
		if err != nil {
			return err
		}
		if cobj == nil {
			return utils.ReplyTextToTelegram(msg, "还不支持捕捉"+target, false)
		}
		objId = cobj.ID
		objName = cobj.Name
		objEmoji = cobj.Emoji
	} else {
		objName = "综合"
		objEmoji = "🏆"
	}

	loots, err := chatLootRank(msg.Chat.ID, objId)
	if err != nil {
		return err
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%s *%s捕捉排行*\n\n", objEmoji, utils.EscapeTelegramMarkdown(objName)))
	if len(loots) == 0 {
		sb.WriteString("无\n")
	} else {
		for i, v := range loots {
			sb.WriteString(fmt.Sprintf("%d\\. %s \\(%d\\)\n", i+1, utils.EscapeTelegramMarkdown(v.Name), v.Amount))
		}
	}
	return utils.ReplyTextToTelegram(msg, sb.String(), true)
}
