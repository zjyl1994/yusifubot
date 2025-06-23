package catchobj

import (
	"errors"

	"github.com/zjyl1994/yusifubot/infra/utils"
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
			return nil, nil
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

func SyncLastObj(db *gorm.DB, user common.UserRel) error {
	var obj CatchObj
	err := db.Where("user_id = ?", user.UserId).Last(&obj).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.NewBizErrWithBase("没有找到任何捕捉配置,请使用setmycatch先创建一个", err)
		}
		return err
	}
	err = CreateCatchObj(db, user, obj.Name, obj.Emoji)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			err = UpdateNickName(db, user, obj.Name)
			if err != nil {
				return err
			}
			err = UpdateEmoji(db, user, obj.Emoji)
			if err != nil {
				return err
			}
			return nil
		}
		return err
	}
	return nil
}
