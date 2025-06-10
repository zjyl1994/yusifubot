package catchobj

import (
	"errors"

	"github.com/zjyl1994/yusifubot/service/catchgame/common"
	"gorm.io/gorm"
)

var ErrCatchObjNotFound = errors.New("catch obj not found")

func ToggleCatch(db *gorm.DB, user common.UserRel) (bool, error) {
	obj, err := GetCatchObj(db, user)
	if err != nil {
		return false, err
	}
	if obj == nil {
		return false, ErrCatchObjNotFound
	}
	obj.Enabled = !obj.Enabled
	return obj.Enabled, db.Save(&obj).Error
}

func GetCatchObj(db *gorm.DB, user common.UserRel) (*CatchObj, error) {
	var obj CatchObj
	err := db.Where(CatchObj{ChatId: user.ChatId, UserId: user.UserId}).First(&obj).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, err
	}
	return &obj, nil
}

func BatchGetCatchObj(db *gorm.DB, chatId int64, userIds []int64) ([]CatchObj, error) {
	var objs []CatchObj
	err := db.Where("chat_id = ? AND user_id IN ?", chatId, userIds).Find(&objs).Error
	if err != nil {
		return nil, err
	}
	return objs, nil
}

func ListEnabledCatchObj(db *gorm.DB, chatId int64) ([]CatchObj, error) {
	var objs []CatchObj
	err := db.Where("chat_id = ? AND enabled = ?", chatId, true).Find(&objs).Error
	if err != nil {
		return nil, err
	}
	return objs, nil
}

func CreateCatchObj(db *gorm.DB, user common.UserRel, name, emoji string) error {
	obj := CatchObj{
		ChatId:  user.ChatId,
		UserId:  user.UserId,
		Enabled: true,
		Name:    name,
		Emoji:   emoji,
	}
	return db.Create(&obj).Error
}

func UpdateNickName(db *gorm.DB, user common.UserRel, name string) error {
	obj, err := GetCatchObj(db, user)
	if err != nil {
		return err
	}
	if obj == nil {
		return ErrCatchObjNotFound
	}
	obj.Name = name
	return db.Save(&obj).Error
}

func UpdateEmoji(db *gorm.DB, user common.UserRel, emoji string) error {
	obj, err := GetCatchObj(db, user)
	if err != nil {
		return err
	}
	if obj == nil {
		return ErrCatchObjNotFound
	}
	obj.Emoji = emoji
	return db.Save(&obj).Error
}
