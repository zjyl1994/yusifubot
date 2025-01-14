package tg

import (
	"errors"
	"strconv"

	"github.com/zjyl1994/yusifubot/infra/vars"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func UpdateUserName(tgUserId int64, tgUserName string) error {
	var m User
	m.UserId = tgUserId
	m.UserName = tgUserName

	return vars.DBInstance.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "tg_user_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"tg_user_name"}),
	}).Create(&m).Error
}

func GetUserName(tgUserId int64) (string, error) {
	var m User
	err := vars.DBInstance.Where("tg_user_id = ?", tgUserId).First(&m).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "USER" + strconv.FormatInt(tgUserId, 10), nil
		}
		return "", err
	}
	return m.UserName, nil
}

func UpdateChatName(tgChatId int64, tgChatName string) error {
	var m Chat
	m.ChatId = tgChatId
	m.ChatName = tgChatName

	return vars.DBInstance.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "tg_chat_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"tg_chat_name"}),
	}).Create(&m).Error
}

func GetChatName(tgChatId int64) (string, error) {
	var m User
	err := vars.DBInstance.Where("tg_chat_id = ?", tgChatId).First(&m).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "CHAT" + strconv.FormatInt(tgChatId, 10), nil
		}
		return "", err
	}
	return m.UserName, nil
}
