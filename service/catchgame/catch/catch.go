package catch

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/zjyl1994/yusifubot/infra/utils"
)

// 结构化后的抓方法
func CatchAction(msg *tgbotapi.Message, catchTarget string, catchNum catchNum) error {
	return utils.ReplyTextToTelegram(msg, fmt.Sprintf("TARGET %s ALL %t NUM %d", catchTarget, catchNum.IsAll(), catchNum.GetNum()), false)
}
