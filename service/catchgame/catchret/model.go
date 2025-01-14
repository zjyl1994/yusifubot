package catchret

type CatchRet struct {
	ID        int64 `gorm:"primaryKey"`
	ChatId    int64 `gorm:"uniqueIndex:idx_chat_user_obj,column:chat_id"`
	UserId    int64 `gorm:"uniqueIndex:idx_chat_user_obj,column:user_id"`
	ObjId     int64 `gorm:"uniqueIndex:idx_chat_user_obj,column:obj_id"`
	Amount    int64 // 捕捉数量
	LastCatch int64 // 最后捕捉时间戳
}

type CatchDetail struct {
	ID        int64 `gorm:"primaryKey"`
	ChatId    int64
	UserId    int64
	ObjId     int64
	Amount    int64
	CatchTime int64
}
