package catch

import (
	"fmt"
	"math/rand/v2"
	"regexp"
	"strings"
	"unicode"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/zjyl1994/yusifubot/infra/utils"
	"github.com/zjyl1994/yusifubot/service/catchgame/catchobj"
	"github.com/zjyl1994/yusifubot/service/catchgame/catchret"
	"github.com/zjyl1994/yusifubot/service/catchgame/common"
	"github.com/zjyl1994/yusifubot/service/catchgame/stamina"
	"github.com/zjyl1994/yusifubot/service/tg"
)

var catchCommandRegexp = regexp.MustCompile(`(?i)^catch([a-zA-Z]+?)?(\d+|all)?$`)

// 所有catch开头命令由此分发
func CatchDispatcher(msg *tgbotapi.Message) error {
	command := msg.Command()
	args := msg.CommandArguments()

	err := tg.UpdateChatAndUserName(msg)
	if err != nil {
		return err
	}
	// 单纯的捕捉指令，合并成组合指令走正则解析
	if command == "catch" {
		var builder strings.Builder
		builder.WriteString(command)
		for _, r := range args {
			if !unicode.IsSpace(r) {
				builder.WriteRune(r)
			}
		}
		command = builder.String()
	}
	// 组合指令，尝试使用正则解析
	if matches := catchCommandRegexp.FindStringSubmatch(command); matches != nil {
		var num catchNum
		catchName := matches[1]
		// catchName为all时代表混抽所有体力
		if strings.EqualFold(catchName, "all") {
			catchName = ""
			num = catchNum("ALL")
		} else { // 明确有抽谁，计算抽的数量
			if len(matches) > 2 && matches[2] != "" {
				num = catchNum(matches[2])
			} else {
				num = catchNum("1")
			}
		}

		return CatchAction(msg, catchName, num)
	}
	return utils.ReplyTextToTelegram(msg, fmt.Sprintf("无法解析：%s %s", command, args), false)
}

// 结构化后的抓方法
func CatchAction(msg *tgbotapi.Message, catchTarget string, catchNum catchNum) (err error) {
	if catchNum.IsAll() || catchNum.GetNum() > 1 { // 抓多个场景单独处理
		return multiCatch(msg, catchTarget, catchNum)
	}
	var cobj *catchobj.CatchObj
	if catchTarget == "" { // 检查是否随机选人抽
		objs, err := catchobj.GetCatchObjs(msg.Chat.ID)
		if err != nil {
			return err
		}
		if len(objs) == 0 {
			return utils.NewBizErr("尚未开放任何捕捉")
		}
		cobj = utils.PickOne(objs) // 随机选择抓取目标
	} else { // 固定抽取
		// 检查抓取对象
		cobj, err = catchobj.GetCatchObjByShorthand(msg.Chat.ID, catchTarget)
		if err != nil {
			return err
		}
		if cobj == nil || (cobj.ChatId != 0 && cobj.ChatId != msg.Chat.ID) || cobj.Stamina == 0 {
			return utils.NewBizErr("尚未开放" + catchTarget + "的捕捉")
		}
	}

	// 消耗用户体力
	user := common.UserRel{
		ChatId: msg.Chat.ID,
		UserId: msg.From.ID,
	}
	_, err = stamina.UseStaminPoint(user, cobj.Stamina)
	if err != nil {
		return err
	}

	if rand.Float64() < cobj.CatchRate { // 计算抓结果
		// 写入结果
		_, err = catchret.AddCatchResult(user, cobj.ID, 1)
		if err != nil {
			return err
		}
		// 生成回复消息
		if len(cobj.CatchHitSticker) > 0 {
			return utils.ReplyStickerToTelegram(msg, cobj.GetHitSticker())
		} else if len(cobj.CatchHitText) > 0 {
			return utils.ReplyTextToTelegram(msg, cobj.GetHitText(), false)
		} else {
			return utils.ReplyTextToTelegram(msg, "成功捕捉"+cobj.Name, false)
		}
	} else {
		if len(cobj.CatchMissSticker) > 0 {
			return utils.ReplyStickerToTelegram(msg, cobj.GetMissSticker())
		} else if len(cobj.CatchMissText) > 0 {
			return utils.ReplyTextToTelegram(msg, cobj.GetMissText(), false)
		} else {
			return utils.ReplyTextToTelegram(msg, "手滑了，"+cobj.Name+"逃走了", false)
		}
	}
}

func multiCatch(msg *tgbotapi.Message, catchTarget string, catchNum catchNum) (err error) {
	// 获取用户体力
	user := common.UserRel{
		ChatId: msg.Chat.ID,
		UserId: msg.From.ID,
	}
	sp, err := stamina.GetStaminPoint(user)
	if err != nil {
		return err
	}
	// 获取需要计算的对象
	var catchObjs []*catchobj.CatchObj
	if catchTarget != "" {
		cobj, err := catchobj.GetCatchObjByShorthand(msg.Chat.ID, catchTarget)
		if err != nil {
			return err
		}
		if cobj == nil || (cobj.ChatId != 0 && cobj.ChatId != msg.Chat.ID) || cobj.Stamina == 0 {
			return utils.NewBizErr("尚未开放" + catchTarget + "的捕捉")
		}
		catchObjs = []*catchobj.CatchObj{cobj}
	} else {
		catchObjs, err = catchobj.GetCatchObjs(msg.Chat.ID)
		if err != nil {
			return err
		}
		if len(catchObjs) == 0 {
			return utils.NewBizErr("尚未开放任何捕捉")
		}
	}
	// 获取抓取列表
	var catchList = make([]*catchobj.CatchObj, 0)
	var remainSp = sp.Current()
	var costSp int64 = 0
	var counter int64 = 0
	for {
		cobj := utils.PickOne(catchObjs)
		if catchNum.IsAll() {
			if cobj.Stamina > remainSp {
				break
			}
		} else {
			if counter >= catchNum.GetNum() {
				break
			}
		}
		remainSp -= cobj.Stamina
		costSp += cobj.Stamina
		counter++
		catchList = append(catchList, cobj)
	}
	if counter == 0 {
		return utils.NewBizErr("体力不足无法捕捉," + sp.String())
	}
	// 消耗用户体力
	_, err = stamina.UseStaminPoint(user, costSp)
	if err != nil {
		return err
	}
	// 计算抓结果
	catchCounterMap := make(map[int64]int64)
	catchNameRel := make(map[int64]string)
	var totalCatch int64
	for i, cobj := range catchList {
		if rand.Float64() < cobj.CatchRate {
			catchNameRel[cobj.ID] = cobj.Name
			catchCounterMap[cobj.ID]++
			totalCatch++
		} else {
			catchList[i] = nil
		}
	}
	// 写入结果
	for cobjID, amount := range catchCounterMap {
		_, err = catchret.AddCatchResult(user, cobjID, amount)
		if err != nil {
			return err
		}
	}
	// 生成回复的消息
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("消耗体力%d, 捕捉%d次，成功率%.2f%%\n", costSp, len(catchList), float64(totalCatch)/float64(len(catchList))*100))
	sb.WriteString("结果：")
	for _, cobj := range catchList {
		if cobj != nil {
			if cobj.Emoji != "" {
				sb.WriteString(cobj.Emoji)
			} else {
				sb.WriteString(CATCH_DEFAULT_EMOJI)
			}
		} else {
			sb.WriteString(CATCH_MISS_EMOJI)
		}
	}
	if totalCatch > 0 {
		sb.WriteString("\n明细如下:\n")
		for cobjID, amount := range catchCounterMap {
			sb.WriteString(fmt.Sprintf("%s:%d\n", catchNameRel[cobjID], amount))
		}
	}
	return utils.ReplyTextToTelegram(msg, sb.String(), false)
}
