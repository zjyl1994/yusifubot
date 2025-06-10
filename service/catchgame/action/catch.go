package action

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/zjyl1994/yusifubot/infra/utils"
	"github.com/zjyl1994/yusifubot/infra/vars"
	"github.com/zjyl1994/yusifubot/service/catchgame/catchobj"
	"github.com/zjyl1994/yusifubot/service/catchgame/common"
	"github.com/zjyl1994/yusifubot/service/catchgame/stamina"
)

func CatchHandler(msg *tgbotapi.Message) (err error) {
	// 获取捕捉对象
	var catchObjList []catchobj.CatchObj
	if msg.ReplyToMessage != nil { // 捕捉特定用户
		cobj, err := catchobj.GetCatchObj(vars.DBInstance, common.UserRel{ChatId: msg.Chat.ID, UserId: msg.ReplyToMessage.From.ID})
		if err != nil {
			return err
		}
		if cobj == nil || !cobj.Enabled {
			return utils.ReplyTextToTelegram(msg, "还不能抓ta", false)
		}
		catchObjList = append(catchObjList, *cobj)
	} else { // 捕捉混池
		catchObjList, err = catchobj.ListEnabledCatchObj(vars.DBInstance, msg.Chat.ID)
		if err != nil {
			return err
		}
	}
	// 计算捕捉数量
	user := common.UserRel{
		ChatId: msg.Chat.ID,
		UserId: msg.From.ID,
	}
	var catchNum int64
	switch msg.Command() {
	case "catch":
		catchNum = 1
	case "catch5":
		catchNum = 5
	case "catch10":
		catchNum = 10
	case "catchall":
		sp, err := stamina.GetStaminPoint(vars.DBInstance, user)
		if err != nil {
			return err
		}
		catchNum = sp.Current() / CATCH_STAMINA_COST
	}
	// 抓
	sp, err := stamina.UseStaminPoint(vars.DBInstance, user, catchNum*CATCH_STAMINA_COST)
	if err != nil {
		return err
	}
	// TODO: finish here
	return nil
}
