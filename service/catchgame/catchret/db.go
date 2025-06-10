package catchret

import (
	"errors"

	"github.com/zjyl1994/yusifubot/infra/utils"
	"github.com/zjyl1994/yusifubot/service/catchgame/common"
	"gorm.io/gorm"
)

func GiveCatchNum(db *gorm.DB, user common.UserRel, objId int64, num int64) error {
	return db.Transaction(func(tx *gorm.DB) error {
		var result CatchResult

		// 根据用户信息捕捉对象查找记录，并使用FOR UPDATE进行锁定
		if err := tx.Set("gorm:query_option", "FOR UPDATE").
			Where("chat_id = ? AND user_id = ? AND obj_id = ?", user.ChatId, user.UserId, objId).First(&result).Error; err != nil {
			// 如果没有找到记录
			if errors.Is(err, gorm.ErrRecordNotFound) {
				if num <= 0 {
					return utils.NewBizErr("现在一只都没有呢~")
				}
				// 创建新的记录
				result = CatchResult{
					ChatId: user.ChatId,
					UserId: user.UserId,
					ObjId:  objId,
					Num:    num,
				}
				return tx.Create(&result).Error
			}
			return err
		}

		// 更新捕捉的数值
		result.Num += num
		if result.Num < 0 {
			result.Num = 0
		}
		// 保存更新后的记录
		return tx.Save(&result).Error
	})
}

func GetMyCatch(db *gorm.DB, user common.UserRel) ([]CatchResult, error) {
	var results []CatchResult
	err := db.Where("chat_id = ? AND user_id = ?", user.ChatId, user.UserId).Find(&results).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []CatchResult{}, nil
		}
		return nil, err
	}
	return results, nil
}

func RankCatch(db *gorm.DB, chatId, objId int64) ([]CatchRankItem, error) {
	var results []CatchRankItem
	query := db.Model(&CatchResult{}).
		Select("user_id, SUM(num) as num").
		Where("chat_id = ?", chatId)

	if objId != 0 {
		query = query.Where("obj_id =?", objId)
	}
	
	err := query.
		Group("user_id").
		Order("num DESC").
		Limit(10).
		Find(&results).Error
	if err != nil {
		return nil, err
	}
	return results, nil
}
