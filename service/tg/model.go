package tg

type User struct {
	ID       int64  `gorm:"primaryKey"`
	UserId   int64  `gorm:"unique,column:tg_user_id"`
	UserName string `gorm:"column:tg_user_name"`
}

type Chat struct {
	ID       int64  `gorm:"primaryKey"`
	ChatId   int64  `gorm:"unique,column:tg_chat_id"`
	ChatName string `gorm:"column:tg_chat_name"`
}
