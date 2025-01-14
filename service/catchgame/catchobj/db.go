package catchobj

import (
	"errors"

	"github.com/zjyl1994/yusifubot/infra/vars"
	"gorm.io/gorm"
)

func GetCatchObj(id int64) (*CatchObj, error) {
	var obj *CatchObj
	err := vars.DBInstance.First(obj, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return obj, nil
}

func GetCatchObjByShorthand(sh string) (*CatchObj, error) {
	var obj *CatchObj
	err := vars.DBInstance.Where("shorthand =?", sh).First(obj).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return obj, nil
}
