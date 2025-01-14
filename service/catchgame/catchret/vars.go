package catchret

import (
	"github.com/zjyl1994/yusifubot/infra/utils/kmutex"
	"github.com/zjyl1994/yusifubot/service/catchgame/common"
)

var retWriteLock = kmutex.NewKmutex(common.UserRelHasher, 100) //体力锁
