package action

import (
	"fmt"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/zjyl1994/yusifubot/infra/utils"
	"github.com/zjyl1994/yusifubot/infra/vars"
	"github.com/zjyl1994/yusifubot/service/catchgame/catchobj"
	"github.com/zjyl1994/yusifubot/service/catchgame/catchret"
	"github.com/zjyl1994/yusifubot/service/catchgame/common"
)

func GetMyCatchHandler(msg *tgbotapi.Message) error {
	user := common.UserRel{ChatId: msg.Chat.ID, UserId: msg.From.ID}
	result, err := catchret.GetMyCatch(vars.DBInstance, user)
	if err != nil {
		return err
	}
	if len(result) == 0 {
		return utils.ReplyTextToTelegram(msg, "两手空空，什么都没抓到", false)
	}
	retMap := make(map[int64]int64)
	objList := make([]int64, 0, len(result))
	for _, v := range result {
		retMap[v.ObjId] = v.Num
		objList = append(objList, v.ObjId)
	}
	if msg.ReplyToMessage != nil { // 回复模式，查询单个对象的
		obj, err := catchobj.GetCatchObj(vars.DBInstance, common.UserRel{ChatId: msg.Chat.ID, UserId: msg.ReplyToMessage.From.ID})
		if err != nil {
			return err
		}
		if obj == nil {
			return utils.ReplyTextToTelegram(msg, "还不能抓ta", false)
		}
		if num, err := retMap[obj.UserId]; err {
			replyText := fmt.Sprintf("已捉到 %s %d只", obj.Name, num)
			return utils.ReplyTextToTelegram(msg, replyText, false)
		} else {
			return utils.ReplyTextToTelegram(msg, "还没捉到ta呢~", false)
		}
	}

	// 直接输入，查询所有
	objs, err := catchobj.BatchGetCatchObj(vars.DBInstance, msg.Chat.ID, objList)
	objNameMap := make(map[int64]string)
	for _, v := range objs {
		objNameMap[v.UserId] = v.Emoji + "" + v.Name
	}
	var sb strings.Builder
	sb.WriteString("**捕获记录**\n\n")
	for _, v := range result {
		name, ok := objNameMap[v.ObjId]
		if !ok {
			name = strconv.FormatInt(v.ObjId, 10)
		}
		sb.WriteString(fmt.Sprintf("%s %d只\n\n", name, v.Num))
	}
	return utils.ReplyTextToTelegram(msg, sb.String(), true)
}
