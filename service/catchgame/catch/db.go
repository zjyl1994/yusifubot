package catch

import (
	"github.com/zjyl1994/yusifubot/infra/vars"
	"github.com/zjyl1994/yusifubot/service/catchgame/common"
)

type lootItem struct {
	Name   string
	Amount int64
}

func lootCount(rel common.UserRel) ([]lootItem, error) {
	querySQL := "SELECT b.`name`,SUM(a.amount) as `amount` FROM catch_rets a LEFT JOIN catch_objs b ON a.obj_id=b.id WHERE a.chat_id = ? AND a.user_id = ?"
	var result []lootItem
	err := vars.DBInstance.Raw(querySQL, rel.ChatId, rel.UserId).Scan(&result).Error
	return result, err
}
