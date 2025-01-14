package catchret

import (
	"time"

	"github.com/zjyl1994/yusifubot/infra/vars"
	"github.com/zjyl1994/yusifubot/service/catchgame/common"
	"gorm.io/gorm"
)

func AddCatchResult(user common.UserRel, objId, amount int64) (int64, error) {
	var detailId int64

	retWriteLock.Lock(user)
	defer retWriteLock.Unlock(user)

	err := vars.DBInstance.Transaction(func(tx *gorm.DB) error {
		now := time.Now().Unix()
		ret := CatchRet{
			ChatId: user.ChatId,
			UserId: user.UserId,
			ObjId:  objId,
		}
		err := tx.FirstOrCreate(&ret).Error
		if err != nil {
			return err
		}

		ret.Amount += amount
		ret.LastCatch = now

		err = tx.Save(ret).Error
		if err != nil {
			return err
		}

		detail := CatchDetail{
			ChatId:    ret.ChatId,
			UserId:    ret.UserId,
			ObjId:     objId,
			Amount:    amount,
			CatchTime: now,
		}
		err = tx.Create(&detail).Error
		if err != nil {
			return err
		}
		detailId = detail.ID
		return nil
	})
	return detailId, err
}
