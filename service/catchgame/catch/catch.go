package catch

import (
	"fmt"
	"math/rand/v2"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
	"github.com/zjyl1994/yusifubot/infra/utils"
	"github.com/zjyl1994/yusifubot/service/catchgame/catchobj"
	"github.com/zjyl1994/yusifubot/service/catchgame/catchret"
	"github.com/zjyl1994/yusifubot/service/catchgame/common"
	"github.com/zjyl1994/yusifubot/service/catchgame/stamina"
	"github.com/zjyl1994/yusifubot/service/tg"
)

var catchCommandRegexp = regexp.MustCompile(`(?i)^catch([a-zA-Z]+?)(\d+|all)?$`)

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
		if len(args) == 0 {
			return utils.ReplyTextToTelegram(msg, "👀 你要捉谁？", false)
		}
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

		if len(matches) > 2 && matches[2] != "" {
			num = catchNum(matches[2])
		} else {
			num = catchNum("1")
		}

		return CatchAction(msg, catchName, num)
	}
	return utils.ReplyTextToTelegram(msg, fmt.Sprintf("无法解析：%s %s", command, args), false)
}

// 结构化后的抓方法
func CatchAction(msg *tgbotapi.Message, catchTarget string, catchNum catchNum) error {
	// 检查抓取对象
	cobj, err := catchobj.GetCatchObjByShorthand(catchTarget)
	if err != nil {
		return err
	}
	if cobj == nil || (cobj.ChatId != 0 && cobj.ChatId != msg.Chat.ID) || cobj.Stamina == 0 {
		return utils.NewBizErr("尚未开放" + catchTarget + "的捕捉")
	}
	logrus.Debugln(cobj)
	// 计算真实抓数
	user := common.UserRel{
		ChatId: msg.Chat.ID,
		UserId: msg.From.ID,
	}
	var realCatchNum int64
	if catchNum.IsAll() {
		sp, err := stamina.GetStaminPoint(user)
		if err != nil {
			return err
		}
		realCatchNum = sp.Current() / cobj.Stamina
		if realCatchNum == 0 {
			realCatchNum = 1
		}
	} else {
		realCatchNum = catchNum.GetNum()
	}
	// 消耗用户体力
	_, err = stamina.UseStaminPoint(user, realCatchNum*cobj.Stamina)
	if err != nil {
		return err
	}
	// 计算抓结果
	catchResult := make([]bool, realCatchNum)
	var catchAmount int64
	for i := range realCatchNum {
		ret := rand.Float64() < cobj.CatchRate
		catchResult[i] = ret
		if ret {
			catchAmount++
		}
	}
	// 写入结果
	_, err = catchret.AddCatchResult(user, cobj.ID, catchAmount)
	if err != nil {
		return err
	}
	// 生成回复的消息
	if realCatchNum == 1 { // 单个捕捉需要支持定制文案和sticker
		if catchResult[0] {
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
	// 多抽模式
	catchSuccessRate := strconv.FormatFloat(float64(catchAmount)/float64(realCatchNum)*100, 'f', 2, 64)
	var sb strings.Builder
	sb.WriteString("捕捉结果：")
	for _, v := range catchResult {
		if v {
			if cobj.Emoji != "" {
				sb.WriteString(cobj.Emoji)
			} else {
				sb.WriteString(CATCH_DEFAULT_EMOJI)
			}
		} else {
			sb.WriteString(CATCH_MISS_EMOJI)
		}
	}
	sb.WriteRune('\n')
	sb.WriteString("本次成功率：")
	sb.WriteString(catchSuccessRate)
	sb.WriteString("%")
	return utils.ReplyTextToTelegram(msg, sb.String(), false)
}
