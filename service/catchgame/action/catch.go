package action

import (
	"fmt"
	"math/rand/v2"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/zjyl1994/yusifubot/infra/utils"
	"github.com/zjyl1994/yusifubot/infra/vars"
	"github.com/zjyl1994/yusifubot/service/catchgame/catchobj"
	"github.com/zjyl1994/yusifubot/service/catchgame/catchret"
	"github.com/zjyl1994/yusifubot/service/catchgame/common"
	"github.com/zjyl1994/yusifubot/service/catchgame/stamina"
)

func CatchHandler(msg *tgbotapi.Message) (err error) {
	user := common.UserRel{
		ChatId: msg.Chat.ID,
		UserId: msg.From.ID,
	}
	catchLock.Lock(user)
	defer catchLock.Unlock(user)
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
	// 消耗体力
	_, err = stamina.UseStaminPoint(vars.DBInstance, user, catchNum*CATCH_STAMINA_COST)
	if err != nil {
		return err
	}
	// 捕捉
	catchResult := make([]catchobj.CatchObj, len(catchObjList))
	catchCount := make(map[catchobj.CatchObj]int64)
	var successCtr int64
	var emojiResult string
	for i := 0; i < len(catchObjList); i++ {
		// 选择捕捉对象
		choiceObj := catchObjList[rand.IntN(len(catchObjList))]
		// 计算是否成功
		success := rand.Float64() < CATCH_RATE
		// 记录成功内容
		if success {
			catchResult[i] = choiceObj
			catchCount[choiceObj]++
			successCtr++

			if choiceObj.Emoji != "" {
				emojiResult += choiceObj.Emoji
			} else {
				emojiResult += CATCH_DEFAULT_EMOJI
			}
		} else {
			emojiResult += CATCH_MISS_EMOJI
		}
	}
	var sb strings.Builder
	sb.WriteString("捕捉结果\n\n")
	sb.WriteString(emojiResult)
	sb.WriteString("\n\n成功率：")
	sb.WriteString(fmt.Sprintf("%.2f%%\n\n\n\n", float64(successCtr)/float64(catchNum)*100))
	for _, obj := range catchResult {
		err = catchret.GiveCatchNum(vars.DBInstance, user, obj.ID, catchCount[obj])
		if err != nil {
			return err
		}
		sb.WriteString(obj.Name)
		sb.WriteString(" ")
		sb.WriteString(strconv.FormatInt(catchCount[obj], 10))
		sb.WriteString("只\n\n")
	}

	return utils.ReplyTextToTelegram(msg, sb.String(), true)
}
