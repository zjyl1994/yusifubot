package action

import (
	"time"

	"github.com/zjyl1994/yusifubot/infra/utils/kmutex"
	"github.com/zjyl1994/yusifubot/service/catchgame/common"
)

const (
	CATCH_MISS_EMOJI    = "✨️" // 未抓到emoji
	CATCH_DEFAULT_EMOJI = "👀"  // 抓到但未设置时的默认emoji
	CATCH_STAMINA_COST  = 10   // 一次抓的SP消耗

	CATCH_AI_JUDGE_KEY      = "catchgame"
	CATCH_AI_JUDGE_COOLDOWN = 10 * time.Second // AI判词冷却时间
	CATCH_AI_CALL_TIMEOUT   = 10 * time.Second // AI最长等待时间
)

type promptData struct {
	UserName string           `json:"玩家"`
	CatchNum int64            `json:"捕捉次数"`
	SuccNum  int64            `json:"成功次数"`
	SuccRate float64          `json:"成功率"`
	Loot     []promptLootItem `json:"捕捉结果"`
}

type promptLootItem struct {
	Name string `json:"对象"`
	Num  int64  `json:"数量"`
}

var (
	catchLock = kmutex.NewKmutex(common.UserRelHasher, 100) //并发捕捉锁
)
