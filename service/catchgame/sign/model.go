package sign

type CatchSign struct {
	ID       int64 `gorm:"primaryKey"`
	ChatId   int64 `gorm:"uniqueIndex:idx_sign_chat_user;column:chat_id"`
	UserId   int64 `gorm:"uniqueIndex:idx_sign_chat_user;column:user_id"`
	SignTime int64
}
