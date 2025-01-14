package tg

type User struct {
	ID       int64 `gorm:"primaryKey"`
	UserId   int64 `gorm:"unique"`
	UserName string
}

func (User) TableName() string {
	return "tg_users"
}

type Chat struct {
	ID       int64 `gorm:"primaryKey"`
	ChatId   int64 `gorm:"unique"`
	ChatName string
}

func (Chat) TableName() string {
	return "tg_chats"
}
