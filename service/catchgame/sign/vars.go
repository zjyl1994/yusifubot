package sign

import (
	"github.com/zjyl1994/yusifubot/infra/utils"
	"github.com/zjyl1994/yusifubot/infra/utils/kmutex"
	"github.com/zjyl1994/yusifubot/service/catchgame/common"
)

var (
	// 默认签到体力范围
	defaultSpMin int64 = 30
	defaultSpMax int64 = 80

	signLock         = kmutex.NewKmutex(common.UserRelHasher, 100)
	ErrAlreadySigned = utils.NewBizErr("今天已经签过到了哦~")
)

type signActivityConfig struct {
	StartTime string `json:"start"`
	EndTime   string `json:"end"`
	Min       int64  `json:"min"`
	Max       int64  `json:"max"`
}
