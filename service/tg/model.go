package tg

type User struct {
	ID       int64 `gorm:"primaryKey"`
	UserId   int64 `gorm:"unique"`
	UserName string
}

type Chat struct {
	ID       int64 `gorm:"primaryKey"`
	ChatId   int64 `gorm:"unique"`
	ChatName string
}
