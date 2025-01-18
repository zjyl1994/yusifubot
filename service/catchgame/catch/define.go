package catch

import (
	"strconv"
	"strings"
)

const (
	CATCH_MISS_EMOJI    = "✨️" // 未抓到emoji
	CATCH_DEFAULT_EMOJI = "👀"  // 抓到但未设置时的默认emoji
)

type catchNum string

func (c catchNum) IsAll() bool {
	return strings.EqualFold(string(c), "all")
}

func (c catchNum) GetNum() int64 {
	num, err := strconv.ParseInt(string(c), 10, 64)
	if err == nil && num > 0 {
		return num
	}
	return 1
}

var statSQLMap = map[string]string{
	"total_rank": `SELECT b.name,SUM(amount) as amount 
	FROM catch_rets a LEFT JOIN catch_objs b ON a.obj_id=b.id
	GROUP BY a.obj_id ORDER BY amount DESC`,
	"chat_rank": `SELECT c.chat_name,b.name,SUM(amount) as amount 
FROM catch_rets a LEFT JOIN catch_objs b ON a.obj_id=b.id
LEFT JOIN tg_chats c ON c.chat_id=a.chat_id
GROUP BY a.chat_id,a.obj_id ORDER BY amount DESC`,
	"user_rank": `SELECT c.user_name,b.name,SUM(amount) as amount 
FROM catch_rets a LEFT JOIN catch_objs b ON a.obj_id=b.id
LEFT JOIN tg_users c ON c.user_id=a.user_id
GROUP BY a.user_id,a.obj_id ORDER BY amount DESC`,
	"day_catch_count": `SELECT date(a.catch_time, 'unixepoch') days,
b.name,SUM(a.amount) amount FROM catch_details a
LEFT JOIN catch_objs b ON a.obj_id=b.id
GROUP BY days ORDER BY days DESC`,
}
