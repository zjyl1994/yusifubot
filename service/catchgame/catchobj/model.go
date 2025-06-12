package catchobj

type CatchObj struct {
	ID      int64 `gorm:"primaryKey"`
	ChatId  int64 `gorm:"uniqueIndex:idx_obj_chat_user;column:chat_id"`
	UserId  int64 `gorm:"uniqueIndex:idx_obj_chat_user;column:user_id"`
	Name    string
	Emoji   string
	Enabled bool
}
