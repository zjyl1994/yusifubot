package action

import (
	"context"
	"fmt"
	"math/rand/v2"
	"strconv"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/zjyl1994/yusifubot/infra/utils"
	"github.com/zjyl1994/yusifubot/infra/vars"
	"github.com/zjyl1994/yusifubot/service/catchgame/catchobj"
	"github.com/zjyl1994/yusifubot/service/catchgame/catchret"
	"github.com/zjyl1994/yusifubot/service/catchgame/common"
	"github.com/zjyl1994/yusifubot/service/catchgame/stamina"
	"github.com/zjyl1994/yusifubot/service/tg"
)

func CatchHandler(msg *models.Message) (err error) {
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
	switch utils.ParseCommand(msg.Text) {
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
		if catchNum <= 0 {
			return utils.ReplyTextToTelegram(msg, "体力不足，"+sp.String(), false)
		}
	}
	// 消耗体力
	_, err = stamina.UseStaminPoint(vars.DBInstance, user, catchNum*CATCH_STAMINA_COST)
	if err != nil {
		return err
	}
	// 捕捉
	catchResult := make([]*catchobj.CatchObj, catchNum)
	catchCount := make(map[catchobj.CatchObj]int64)
	randResult := make([]byte, catchNum)
	_, err = vars.RNG.Read(randResult)
	if err != nil {
		return err
	}

	var successCtr int64
	var emojiResult string
	for i := range catchNum { // 多轮捕捉
		// 选择捕捉对象
		choiceObj := catchObjList[rand.N(len(catchObjList))]
		// 计算是否成功
		success := randResult[i] > 128
		// 记录成功内容
		if success {
			catchResult[i] = &choiceObj
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
	sb.WriteString("<b>捕捉结果</b>\n")
	sb.WriteString(emojiResult)
	if successCtr == 0 {
		sb.WriteRune('\n')
		sb.WriteString(tg.GetTgUserName(msg.From))
		sb.WriteString("两手空空不知所措")
	} else {
		sb.WriteString("<blockquote expandable>")
		succRate := float64(successCtr) / float64(catchNum)
		sb.WriteString(fmt.Sprintf("成功率： %.2f%%\n\n", succRate*100))
		for obj, num := range catchCount {
			err = catchret.GiveCatchNum(vars.DBInstance, user, obj.UserId, num)
			if err != nil {
				return err
			}
			sb.WriteString(obj.Name)
			sb.WriteString(" ")
			sb.WriteString(strconv.FormatInt(num, 10))
			sb.WriteString("只\n")
		}
		sb.WriteString("</blockquote>")
	}
	// 回复给玩家消息
	var msgParams bot.SendMessageParams
	msgParams.Text = sb.String()
	msgParams.ReplyParameters.MessageID = msg.ID
	msgParams.ParseMode = models.ParseModeHTML
	msgParams.MessageEffectID = "5046509860389126442"
	_, err = vars.BotInstance.SendMessage(context.Background(), &msgParams)
	return err
}
