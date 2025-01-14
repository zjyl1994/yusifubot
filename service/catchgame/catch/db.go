package catch

import (
	"github.com/zjyl1994/yusifubot/infra/vars"
	"github.com/zjyl1994/yusifubot/service/catchgame/common"
)

type lootItem struct {
	Name   string
	Amount int64
}

func userLootCount(rel common.UserRel) ([]lootItem, error) {
	querySQL := "SELECT b.`name`,SUM(a.amount) as `amount` FROM catch_rets a LEFT JOIN catch_objs b ON a.obj_id=b.id WHERE a.chat_id = ? AND a.user_id = ? ORDER BY `amount` DESC"
	var result []lootItem
	err := vars.DBInstance.Raw(querySQL, rel.ChatId, rel.UserId).Scan(&result).Error
	return result, err
}

func chatLootRank(chatId, objId int64) ([]lootItem, error) {
	query := vars.DBInstance
	if objId == 0 { // 全局排行
		sql := "SELECT b.user_name AS `name`,SUM(a.amount) AS `amount` FROM catch_rets a LEFT JOIN tg_users b ON a.user_id = b.user_id WHERE a.chat_id = ?  ORDER BY `amount` DESC LIMIT 10"
		query = query.Raw(sql, chatId)
	} else { // 分类排行
		sql := "SELECT b.user_name AS `name`,SUM(a.amount) AS `amount` FROM catch_rets a LEFT JOIN tg_users b ON a.user_id = b.user_id WHERE a.chat_id = ? AND a.obj_id = ?  ORDER BY `amount` DESC LIMIT 10"
		query = query.Raw(sql, chatId, objId)
	}

	var result []lootItem
	err := query.Scan(&result).Error
	return result, err
}
