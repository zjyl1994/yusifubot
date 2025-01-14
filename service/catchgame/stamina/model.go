package stamina

import (
	"fmt"
	"time"

	"github.com/zjyl1994/yusifubot/infra/utils"
)

type Stamina struct {
	ID       int64 `gorm:"primaryKey"`
	ChatId   int64 `gorm:"uniqueIndex:idx_chat_user;column:chat_id"`
	UserId   int64 `gorm:"uniqueIndex:idx_chat_user;column:user_id"`
	LastTick int64
	LastSP   int64
}

// 计算当前体力
func (s *Stamina) Current() int64 {
	return utils.IdleCalcWithMax(s.LastTick, time.Now().Unix(), s.LastSP, SP_PER_SECOND, SP_MAX)
}

// 计算恢复下一点体力的剩余秒数
func (s *Stamina) RemainSecond() int64 {
	current := s.Current()
	if current >= SP_MAX { // 体力满不会自动回复
		return 0
	}
	elapsedSecond := time.Now().Unix() - s.LastTick
	return SP_PER_SECOND - (elapsedSecond % SP_PER_SECOND)
}

func (s *Stamina) String() string {
	current := s.Current()
	remainSec := s.RemainSecond()
	if remainSec > 0 {
		return fmt.Sprintf("当前体力剩余%d,距离恢复下一点还有%d秒。", current, remainSec)
	} else {
		return fmt.Sprintf("当前体力剩余%d,已达到自然恢复上限。", current)
	}
}
