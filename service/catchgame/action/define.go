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
)

var (
	catchLock = kmutex.NewKmutex(common.UserRelHasher, 100) //并发捕捉锁
)
