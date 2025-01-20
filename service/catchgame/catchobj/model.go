package catchobj

import (
	"github.com/zjyl1994/yusifubot/infra/utils"
)

type CatchObj struct {
	ID               int64               `gorm:"primaryKey"`
	ChatId           int64               `gorm:"uniqueIndex:idx_chat_short;column:chat_id"`
	Shorthand        string              `gorm:"uniqueIndex:idx_chat_short;column:shorthand"` // 缩写
	Name             string              // 捕捉的人
	Emoji            string              // 缩写的emoji
	Stamina          int64               // 捕捉消耗的体力
	CatchRate        float64             // 爆率
	CatchMissText    utils.Array[string] `gorm:"type:TEXT;serializer:json"` // 捕捉失败提示语
	CatchMissSticker utils.Array[string] `gorm:"type:TEXT;serializer:json"` // 捕捉失败贴图
	CatchHitText     utils.Array[string] `gorm:"type:TEXT;serializer:json"` // 捕捉成功提示语
	CatchHitSticker  utils.Array[string] `gorm:"type:TEXT;serializer:json"` // 捕捉成功贴图
	TgUserId         string              // Tg用户id
}

func (obj CatchObj) GetMissText() string {
	return utils.PickOne(obj.CatchMissText)
}

func (obj CatchObj) GetMissSticker() string {
	return utils.PickOne(obj.CatchMissSticker)
}

func (obj CatchObj) GetHitText() string {
	return utils.PickOne(obj.CatchHitText)
}

func (obj CatchObj) GetHitSticker() string {
	return utils.PickOne(obj.CatchHitSticker)
}
