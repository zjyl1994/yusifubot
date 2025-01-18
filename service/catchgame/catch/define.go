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
	"全服捕捉榜": `SELECT b.name AS 对象,SUM(amount) as 数量 
	FROM catch_rets a LEFT JOIN catch_objs b ON a.obj_id=b.id
	GROUP BY a.obj_id ORDER BY 数量 DESC`,
	"会话捕捉榜": `SELECT c.chat_name AS 会话,b.name AS 对象,SUM(amount) as 数量 
FROM catch_rets a LEFT JOIN catch_objs b ON a.obj_id=b.id
LEFT JOIN tg_chats c ON c.chat_id=a.chat_id
GROUP BY a.chat_id,a.obj_id ORDER BY 数量 DESC`,
	"用户捕捉榜": `SELECT c.user_name AS 用户,b.name AS 对象,SUM(amount) as 数量 
FROM catch_rets a LEFT JOIN catch_objs b ON a.obj_id=b.id
LEFT JOIN tg_users c ON c.user_id=a.user_id
GROUP BY a.user_id,a.obj_id ORDER BY 数量 DESC`,
	"日捕捉榜": `SELECT date(a.catch_time, 'unixepoch') 日期,
b.name AS 对象,SUM(a.amount) 数量 FROM catch_details a
LEFT JOIN catch_objs b ON a.obj_id=b.id
GROUP BY 日期 ORDER BY 日期 DESC`,
}
