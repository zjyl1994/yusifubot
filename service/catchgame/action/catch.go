package action

import (
	"context"
	"fmt"
	"math/rand/v2"
	"strconv"
	"strings"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/sirupsen/logrus"
	"github.com/zjyl1994/yusifubot/infra/utils"
	"github.com/zjyl1994/yusifubot/infra/vars"
	"github.com/zjyl1994/yusifubot/service/catchgame/catchobj"
	"github.com/zjyl1994/yusifubot/service/catchgame/catchret"
	"github.com/zjyl1994/yusifubot/service/catchgame/common"
	"github.com/zjyl1994/yusifubot/service/catchgame/stamina"
	"github.com/zjyl1994/yusifubot/service/config"
	"github.com/zjyl1994/yusifubot/service/tg"
)

func CatchHandler(msg *models.Message) (err error) {
	// 只在群聊中生效
	if !utils.IsGroup(msg) {
		return nil
	}

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
	if len(catchObjList) == 0 {
		return utils.ReplyTextToTelegram(msg, "没有人可以捕捉哦", false)
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
	succRate := (float64(successCtr) / float64(catchNum)) * 100
	// 生成捕捉结果
	// 1. 判断对话内是否开启AI生成评判词
	var aiEnabled bool
	if val, err := config.Get(msg.Chat.ID, "ai"); err == nil {
		aiEnabled, err = strconv.ParseBool(val)
		if err != nil {
			logrus.Errorf("catch game parse ai config failed: %v", err)
		}
	} else {
		logrus.Errorf("catch game get ai config failed: %v", err)
	}
	// 生成AI判词
	var aiJudgment string
	if successCtr > 0 && aiEnabled && vars.ReplicateCooldown.CheckAndSetCooldown(CATCH_AI_JUDGE_KEY, CATCH_AI_JUDGE_COOLDOWN) {
		const (
			SYSTEM_PROMPT         = "群中正在进行一场捕捉游戏，你作为一位旁观者对捕捉结果进行简单评论。评价结果请用中文回复。"
			MAX_COMPLETION_TOKENS = 256
			TEMPERATURE           = 1.1
		)
		prompt := fmt.Sprintf("玩家'%s'在本轮捕捉中成功率%.2f,战利品如下:\n", tg.GetTgUserName(msg.From), succRate)
		for obj, num := range catchCount {
			prompt += fmt.Sprintf("%s: %d只\n", obj.Name, num)
		}
		gptResp, err := utils.CallGPT41Nano(SYSTEM_PROMPT, prompt, TEMPERATURE, MAX_COMPLETION_TOKENS, 10*time.Second)
		if err != nil {
			logrus.Errorf("catch game replicate request failed: %v", err)
		} else {
			aiJudgment = gptResp
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
		if aiEnabled && len(aiJudgment) > 0 {
			sb.WriteString(aiJudgment)
		} else {
			succRate := float64(successCtr) / float64(catchNum)
			sb.WriteString(fmt.Sprintf("成功率： %.2f%%\n", succRate*100))
		}
		sb.WriteString("\n本轮战利品：\n")
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
	msgParams.ChatID = msg.Chat.ID
	msgParams.ReplyParameters = &models.ReplyParameters{
		MessageID: msg.ID,
	}
	msgParams.DisableNotification = true
	msgParams.ParseMode = models.ParseModeHTML
	_, err = vars.BotInstance.SendMessage(context.Background(), &msgParams)
	return err
}
