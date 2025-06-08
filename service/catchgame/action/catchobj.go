package action

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/zjyl1994/yusifubot/infra/utils"
	"github.com/zjyl1994/yusifubot/service/catchgame/catchobj"
	"github.com/zjyl1994/yusifubot/service/catchgame/common"
	"github.com/zjyl1994/yusifubot/service/tg"

	"regexp"
)

// 切换可抓状态
func CatchMeHandler(msg *tgbotapi.Message) error {
	current, err := catchobj.ToggleCatch(common.UserRel{ChatId: msg.Chat.ID, UserId: msg.From.ID})
	if err != nil {
		if err == catchobj.ErrCatchObjNotFound {
			err = createCatchObj(msg, "", "")
			if err != nil {
				return err
			}
			return utils.ReplyTextToTelegram(msg, "成功开启捕捉自己功能\n\n可以使用setnickname自定义昵称\n\n可以使用setemoji自定义emoji", true)
		}
		return err
	}
	var reply string
	if current {
		reply = "成功开启群内的捕捉自己功能"
	} else {
		reply = "您关闭了群内的捕捉自己功能"
	}
	return utils.ReplyTextToTelegram(msg, reply, false)
}

// 设置昵称
func SetMyNicknameHandler(msg *tgbotapi.Message) error {
	arg := msg.CommandArguments()
	err := catchobj.UpdateNickName(common.UserRel{ChatId: msg.Chat.ID, UserId: msg.From.ID}, arg)
	if err == catchobj.ErrCatchObjNotFound {
		err = createCatchObj(msg, arg, "")
	}
	if err != nil {
		return err
	}
	return utils.ReplyTextToTelegram(msg, "成功设置昵称为"+arg, false)
}

// 设置emoji
func SetMyEmojiHandler(msg *tgbotapi.Message) error {
	arg := msg.CommandArguments()
	if !isSingleEmoji(arg) {
		return utils.ReplyTextToTelegram(msg, "只能用一个emoji哦", false)
	}
	err := catchobj.UpdateEmoji(common.UserRel{ChatId: msg.Chat.ID, UserId: msg.From.ID}, arg)
	if err == catchobj.ErrCatchObjNotFound {
		err = createCatchObj(msg, "", arg)
	}
	if err != nil {
		return err
	}
	return utils.ReplyTextToTelegram(msg, "成功设置emoji为"+arg, false)
}

// 静默创建可抓账号
func createCatchObj(msg *tgbotapi.Message, nickName, emoji string) error {
	if nickName == "" {
		nickName = tg.GetTgUserName(msg.From)
	}
	if emoji == "" {
		emoji = CATCH_DEFAULT_EMOJI
	}
	return catchobj.CreateCatchObj(common.UserRel{ChatId: msg.Chat.ID, UserId: msg.From.ID}, nickName, emoji)
}

// 判断字符串是否是一个单一的 Emoji
func isSingleEmoji(s string) bool {
	// 正则表达式匹配一个完整的 Emoji（包括复合 Emoji）
	// 来源：简化版 emoji 正则，适用于常见场景
	re := regexp.MustCompile(`^([\p{Emoji}\p{Emoticons}][\uFE00-\uFE0F]?|[\p{Emoji}\p{Emoticons}]\u200D[\p{Emoji}\p{Emoticons}])+$`)
	return re.MatchString(s)
}
