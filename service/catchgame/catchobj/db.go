package catchobj

import (
	"errors"

	"github.com/zjyl1994/yusifubot/infra/vars"
	"github.com/zjyl1994/yusifubot/service/catchgame/common"
	"gorm.io/gorm"
)

var ErrCatchObjNotFound = errors.New("catch obj not found")

func ToggleCatch(user common.UserRel) (bool, error) {
	obj, err := GetCatchObj(user)
	if err != nil {
		return false, err
	}
	if obj == nil {
		return false, ErrCatchObjNotFound
	}
	obj.Enabled = !obj.Enabled
	return obj.Enabled, vars.DBInstance.Save(&obj).Error
}

func GetCatchObj(user common.UserRel) (*CatchObj, error) {
	var obj CatchObj
	err := vars.DBInstance.Where(CatchObj{ChatId: user.ChatId, UserId: user.UserId}).First(&obj).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, err
	}
	return &obj, nil
}

func CreateCatchObj(user common.UserRel, name, emoji string) error {
	obj := CatchObj{
		ChatId:  user.ChatId,
		UserId:  user.UserId,
		Enabled: true,
		Name:    name,
		Emoji:   emoji,
	}
	return vars.DBInstance.Create(&obj).Error
}

func UpdateNickName(user common.UserRel, name string) error {
	obj, err := GetCatchObj(user)
	if err != nil {
		return err
	}
	if obj == nil {
		return ErrCatchObjNotFound
	}
	obj.Name = name
	return vars.DBInstance.Save(&obj).Error
}

func UpdateEmoji(user common.UserRel, emoji string) error {
	obj, err := GetCatchObj(user)
	if err != nil {
		return err
	}
	if obj == nil {
		return ErrCatchObjNotFound
	}
	obj.Emoji = emoji
	return vars.DBInstance.Save(&obj).Error
}
