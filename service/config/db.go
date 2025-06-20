package config

import (
	"errors"

	"github.com/zjyl1994/yusifubot/infra/vars"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func Get(chatId int64, name string) (string, error) {
	var m ChatConfig
	err := vars.DBInstance.Where("chat_id = ? AND config_key = ?", chatId, name).Find(&m).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", nil
		}
		return "", err
	}
	return m.ConfigData, nil
}

func Set(chatId int64, name, data string) error {
	m := ChatConfig{
		ChatId:     chatId,
		ConfigKey:  name,
		ConfigData: data,
	}
	return vars.DBInstance.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "chat_id"}, {Name: "config_key"}},
		DoUpdates: clause.AssignmentColumns([]string{"config_data"}),
	}).Create(&m).Error
}
