package draw

import (
	"sync"

	"github.com/go-telegram/bot/models"
	"github.com/zjyl1994/yusifubot/infra/utils"
	"github.com/zjyl1994/yusifubot/infra/vars"
)

var allowChatIds = make(map[int64]struct{})
var allowChatMutex sync.Mutex

func SwitchHandler(msg *models.Message) error {
	if msg.From.ID != vars.AdminUserId {
		return utils.ReplyTextToTelegram(msg, "您不是管理员", false)
	}

	allowed := toggleChatAllowance(msg.Chat.ID)
	if allowed {
		return utils.ReplyTextToTelegram(msg, "已开启绘图能力", false)
	} else {
		return utils.ReplyTextToTelegram(msg, "已关闭绘图能力", false)
	}
}

func toggleChatAllowance(chatId int64) bool {
	allowChatMutex.Lock()
	defer allowChatMutex.Unlock()

	if _, ok := allowChatIds[chatId]; ok {
		delete(allowChatIds, chatId)
		return false
	} else {
		allowChatIds[chatId] = struct{}{}
		return true
	}
}

func checkChatIdAllowed(chatId int64) bool {
	allowChatMutex.Lock()
	defer allowChatMutex.Unlock()

	_, ok := allowChatIds[chatId]
	return ok
}
