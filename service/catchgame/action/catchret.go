package action

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/go-telegram/bot/models"
	"github.com/zjyl1994/yusifubot/infra/utils"
	"github.com/zjyl1994/yusifubot/infra/vars"
	"github.com/zjyl1994/yusifubot/service/catchgame/catchobj"
	"github.com/zjyl1994/yusifubot/service/catchgame/catchret"
	"github.com/zjyl1994/yusifubot/service/catchgame/common"
	"github.com/zjyl1994/yusifubot/service/tg"
)

func GetMyCatchHandler(msg *models.Message) error {
	// 只在群聊中生效
	if !utils.IsGroup(msg) {
		return nil
	}

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

func RankCatchHandler(msg *models.Message) error {
	// 只在群聊中生效
	if !utils.IsGroup(msg) {
		return nil
	}

	var objId int64
	if msg.ReplyToMessage != nil {
		objId = msg.ReplyToMessage.From.ID
	}
	items, err := catchret.RankCatch(vars.DBInstance, msg.Chat.ID, objId)
	if err != nil {
		return err
	}
	if len(items) == 0 {
		return utils.ReplyTextToTelegram(msg, "现在还没有人捉到", false)
	}
	var sb strings.Builder
	if objId != 0 {
		obj, err := catchobj.GetCatchObj(vars.DBInstance, common.UserRel{ChatId: msg.Chat.ID, UserId: objId})
		if err != nil {
			return err
		}
		if obj == nil {
			return utils.ReplyTextToTelegram(msg, "还不能抓ta", false)
		}
		sb.WriteString(fmt.Sprintf("**%s 捕捉排行榜**\n\n", obj.Name))
	} else {
		sb.WriteString("**🏆综合捕捉排行榜**\n\n")
	}
	for idx, v := range items {
		name, err := tg.GetUserName(v.UserId)
		if err != nil {
			return err
		}
		sb.WriteString(fmt.Sprintf("%d. %s %d只\n\n", idx+1, name, v.Num))
	}
	return utils.ReplyTextToTelegram(msg, sb.String(), true)
}
