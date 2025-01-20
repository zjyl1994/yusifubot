package catchobj

import (
	"errors"

	"github.com/zjyl1994/yusifubot/infra/vars"
	"gorm.io/gorm"
)

func GetCatchObj(id int64) (*CatchObj, error) {
	var obj CatchObj
	err := vars.DBInstance.First(&obj, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &obj, nil
}

func GetCatchObjByShorthand(chatId int64, sh string) (*CatchObj, error) {
	var obj CatchObj
	err := vars.DBInstance.Where("(chat_id = ? OR chat_id = 0)", chatId).Where("shorthand =?", sh).First(&obj).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &obj, nil
}

func GetCatchObjs(chatId int64) ([]*CatchObj, error) {
	var objs []*CatchObj
	err := vars.DBInstance.Where("(chat_id = ? OR chat_id = 0)", chatId).Find(&objs).Error
	if err != nil {
		return nil, err
	}
	return objs, nil
}
