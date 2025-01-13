package stamina

import (
	"errors"
	"strconv"

	"github.com/zjyl1994/yusifubot/infra/utils/kmutex"
	"github.com/zjyl1994/yusifubot/service/catchgame/common"
)

const (
	SP_MAX        = 200 // 体力自然恢复上限
	SP_PER_SECOND = 400 // 400秒一点体力
)

var (
	spLock       = kmutex.NewKmutex(userRelHasher, 0) //体力锁
	ErrNotEnough = errors.New("体力不足")
)

func userRelHasher(rel common.UserRel) uint64 {
	return kmutex.StringHasher(strconv.FormatInt(rel.ChatId, 10) + "_" + strconv.FormatInt(rel.UserId, 10))
}
