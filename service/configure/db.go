package configure

import (
	"errors"

	"github.com/zjyl1994/yusifubot/infra/utils"
	"github.com/zjyl1994/yusifubot/infra/vars"
	"golang.org/x/sync/singleflight"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	cacheMap = utils.NewSafeMap[string, string]()
	cacheSf  singleflight.Group
)

func Get(name, defaultValue string) (string, error) {
	val, ok := cacheMap.Get(name)
	if ok {
		return val, nil
	}
	ret, err, _ := cacheSf.Do(name, func() (any, error) {
		var m Configure
		err := vars.DBInstance.Where("name = ?", name).First(&m).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				cacheMap.Set(name, defaultValue)
				return defaultValue, nil
			}
			return "", err
		}
		cacheMap.Set(name, m.Data)
		return m.Data, nil
	})
	if err != nil {
		return "", err
	}
	return ret.(string), nil
}

func Set(name, value string) error {
	m := &Configure{
		Name: name,
		Data: value,
	}
	err := vars.DBInstance.Clauses(clause.OnConflict{UpdateAll: true}).Create(m).Error
	if err != nil {
		return err
	}

	cacheMap.Set(name, value)
	return nil
}
