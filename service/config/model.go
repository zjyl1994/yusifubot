package config

type ChatConfig struct {
	ID         int64  `gorm:"primaryKey"`
	ChatId     int64  `gorm:"uniqueIndex:idx_cc_chat_key;column:chat_id"`
	ConfigKey  string `gorm:"uniqueIndex:idx_cc_chat_key;column:config_key"`
	ConfigData string
}
