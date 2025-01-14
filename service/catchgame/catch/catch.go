package catch

import (
	"math/rand/v2"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
	"github.com/zjyl1994/yusifubot/infra/utils"
	"github.com/zjyl1994/yusifubot/service/catchgame/catchobj"
	"github.com/zjyl1994/yusifubot/service/catchgame/catchret"
	"github.com/zjyl1994/yusifubot/service/catchgame/common"
	"github.com/zjyl1994/yusifubot/service/catchgame/stamina"
)

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
