package catchret

type CatchResult struct {
	ID     int64 `gorm:"primaryKey"`
	ChatId int64 `gorm:"uniqueIndex:idx_chat_user_obj;column:chat_id"`
	UserId int64 `gorm:"uniqueIndex:idx_chat_user_obj;column:user_id"`
	ObjId  int64 `gorm:"uniqueIndex:idx_chat_user_obj;column:obj_id"`
	Num    int64
}
