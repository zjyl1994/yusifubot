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

type statSQLItem struct {
	Name string
	Desc string
	SQL  string
}

type StatResult struct {
	Name   string
	Desc   string
	Result []map[string]any
}

var statSQLs = []statSQLItem{
	{
		Name: "捕捉倍率公示",
		Desc: "为保障公平公正公开，特在此公示登记对象的捕捉概率",
		SQL: `SELECT
  IFNULL( b.chat_name, "(全服)" ) AS 会话,
  a.name AS 对象,
  a.stamina AS 消耗体力,
  printf ( "%.2f%%", a.catch_rate * 100 ) AS 爆率 
FROM
  catch_objs a
  LEFT JOIN tg_chats b ON a.chat_id = b.chat_id`,
	},
	{
		Name: "今日捕捉榜",
		Desc: "最多展示前30条记录",
		SQL: `SELECT b.name AS 对象,
SUM(a.amount) 数量 FROM catch_details a
LEFT JOIN catch_objs b ON a.obj_id=b.id
WHERE a.catch_time >= strftime('%s', date('now', 'start of day'))
GROUP BY a.obj_id HAVING 数量 > 0 ORDER BY 数量 DESC LIMIT 30`,
	},
	{
		Name: "全服捕捉榜",
		Desc: "最多展示前30条记录",
		SQL: `SELECT b.name AS 对象,SUM(amount) as 数量 
	FROM catch_rets a LEFT JOIN catch_objs b ON a.obj_id=b.id
	GROUP BY a.obj_id HAVING 数量 > 0 ORDER BY 数量 DESC LIMIT 30`,
	},
	{
		Name: "会话捕捉榜",
		Desc: "最多展示前30条记录",
		SQL: `SELECT c.chat_name AS 会话,b.name AS 对象,SUM(amount) as 数量 
FROM catch_rets a LEFT JOIN catch_objs b ON a.obj_id=b.id
LEFT JOIN tg_chats c ON c.chat_id=a.chat_id
GROUP BY a.chat_id,a.obj_id HAVING 数量 > 0 ORDER BY 数量 DESC LIMIT 30`,
	},
	{
		Name: "用户捕捉榜",
		Desc: "最多展示前30条记录",
		SQL: `SELECT c.user_name AS 用户,b.name AS 对象,SUM(amount) as 数量 
FROM catch_rets a LEFT JOIN catch_objs b ON a.obj_id=b.id
LEFT JOIN tg_users c ON c.user_id=a.user_id
GROUP BY a.user_id,a.obj_id HAVING 数量 > 0 ORDER BY 数量 DESC LIMIT 30`,
	},
}
