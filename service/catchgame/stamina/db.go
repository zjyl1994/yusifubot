package stamina

import (
	"errors"
	"fmt"
	"time"

	"github.com/zjyl1994/yusifubot/infra/utils"
	"github.com/zjyl1994/yusifubot/service/catchgame/common"
	"gorm.io/gorm"
)

func GetStaminPoint(db *gorm.DB, user common.UserRel) (*Stamina, error) {
	var sp Stamina
	err := db.Where(Stamina{ChatId: user.ChatId, UserId: user.UserId}).First(&sp).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) { // 不存在记录的用户直接返回最高能量上限
			sp.ChatId = user.ChatId
			sp.UserId = user.UserId
			sp.LastSP = SP_MAX
			sp.LastTick = time.Now().Unix()
			return &sp, nil
		}
		return nil, err
	}
	return &sp, nil
}

func UseStaminPoint(db *gorm.DB, user common.UserRel, cost int64) (*Stamina, error) {
	spLock.Lock(user)
	defer spLock.Unlock(user)
	
	var sp *Stamina
	err := db.Transaction(func(tx *gorm.DB) error {
		var err error
		sp, err = GetStaminPoint(tx, user)
		if err != nil {
			return err
		}
		// 计算并扣减能量
		current := sp.Current()
		remainEnergy := current - cost
		// 检查是否扣完
		if remainEnergy < 0 {
			return utils.NewBizErrWithBase(fmt.Sprintf("体力不足%d,%s", cost, sp.String()), ErrNotEnough)
		}
		// 新体力写入DB
		sp.LastSP = remainEnergy
		sp.LastTick = time.Now().Unix()

		return tx.Save(sp).Error
	})
	
	if err != nil {
		return nil, err
	}
	return sp, nil
}

func AddStaminPoint(db *gorm.DB, user common.UserRel, amount int64) error {
	spLock.Lock(user)
	defer spLock.Unlock(user)
	
	return db.Transaction(func(tx *gorm.DB) error {
		sp, err := GetStaminPoint(tx, user)
		if err != nil {
			return err
		}

		sp.LastSP = sp.Current() + amount
		sp.LastTick = time.Now().Unix()

		return tx.Save(sp).Error
	})
}
