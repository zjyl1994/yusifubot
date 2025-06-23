package action

import (
	"strconv"
	"unicode/utf8"

	"github.com/go-telegram/bot/models"
	"github.com/zjyl1994/yusifubot/infra/utils"
	"github.com/zjyl1994/yusifubot/infra/vars"
	"github.com/zjyl1994/yusifubot/service/catchgame/catchobj"
	"github.com/zjyl1994/yusifubot/service/catchgame/catchret"
	"github.com/zjyl1994/yusifubot/service/catchgame/common"
	"github.com/zjyl1994/yusifubot/service/tg"
	"golang.org/x/text/unicode/norm"
)

// 切换可抓状态
func CatchMeHandler(msg *models.Message) error {
	// 只在群聊中生效
	if !utils.IsGroup(msg) {
		return nil
	}

	if msg.ReplyToMessage != nil { // 送自己给某人
		sendUser := msg.ReplyToMessage.From
		args := utils.ParseCommandArguments(msg.Text)
		if len(args) == 0 {
			return utils.ReplyTextToTelegram(msg, "请在命令后追加要赠送的数量", false)
		}
		num, err := strconv.Atoi(args[0])
		if err != nil {
			return utils.ReplyTextToTelegram(msg, "请在命令后追加要赠送的数量", false)
		}
		err = catchret.GiveCatchNum(vars.DBInstance, common.UserRel{ChatId: msg.Chat.ID, UserId: sendUser.ID}, msg.From.ID, int64(num))
		if err != nil {
			return err
		}
		return utils.ReplyTextToTelegram(msg, "成功赠送"+strconv.Itoa(num)+"只", false)
	}

	current, err := catchobj.ToggleCatch(vars.DBInstance, common.UserRel{ChatId: msg.Chat.ID, UserId: msg.From.ID})
	if err != nil {
		if err == catchobj.ErrCatchObjNotFound {
			err = createCatchObj(msg, "", "")
			if err != nil {
				return err
			}
			return utils.ReplyTextToTelegram(msg, "成功开启捕捉自己功能\n\n可以使用setmyname自定义昵称\n\n可以使用setmyemoji自定义emoji", true)
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

func SetMyCatchHandler(msg *models.Message) error {
	// 只在群聊中生效
	if !utils.IsGroup(msg) {
		return nil
	}
	args := utils.ParseCommandArguments(msg.Text)
	switch len(args) {
	case 0:
		return utils.ReplyTextToTelegram(msg, "请在命令后追加要设置的新昵称和emoji", false)
	case 1:
		if isSingleEmoji(args[0]) {
			err := catchobj.UpdateEmoji(vars.DBInstance, common.UserRel{ChatId: msg.Chat.ID, UserId: msg.From.ID}, args[0])
			if err == catchobj.ErrCatchObjNotFound {
				err = createCatchObj(msg, tg.GetTgUserName(msg.From), args[0])
			}
			if err != nil {
				return err
			}
			return utils.ReplyTextToTelegram(msg, "成功设置emoji", false)
		} else {
			err := catchobj.UpdateNickName(vars.DBInstance, common.UserRel{ChatId: msg.Chat.ID, UserId: msg.From.ID}, args[0])
			if err == catchobj.ErrCatchObjNotFound {
				err = createCatchObj(msg, args[0], CATCH_DEFAULT_EMOJI)
			}
			if err != nil {
				return err
			}
			return utils.ReplyTextToTelegram(msg, "成功设置昵称", false)
		}
	default:
		catchName := args[0]
		catchEmoji := args[1]
		if !isSingleEmoji(catchEmoji) {
			return utils.ReplyTextToTelegram(msg, "只能用一个emoji哦", false)
		}
		cobj, err := catchobj.GetCatchObj(vars.DBInstance, common.UserRel{ChatId: msg.Chat.ID, UserId: msg.From.ID})
		if err != nil {
			return err
		}
		if cobj == nil {
			err = createCatchObj(msg, catchName, catchEmoji)
			if err != nil {
				return err
			}
			return utils.ReplyTextToTelegram(msg, "成功设置昵称和emoji", false)
		}
		err = catchobj.UpdateNickName(vars.DBInstance, common.UserRel{ChatId: msg.Chat.ID, UserId: msg.From.ID}, catchName)
		if err != nil {
			return err
		}
		err = catchobj.UpdateEmoji(vars.DBInstance, common.UserRel{ChatId: msg.Chat.ID, UserId: msg.From.ID}, catchEmoji)
		if err != nil {
			return err
		}
		return utils.ReplyTextToTelegram(msg, "成功设置昵称和emoji", false)
	}
}

// 设置昵称
func SetMyNicknameHandler(msg *models.Message) error {
	// 只在群聊中生效
	if !utils.IsGroup(msg) {
		return nil
	}

	args := utils.ParseCommandArguments(msg.Text)
	if len(args) == 0 {
		return utils.ReplyTextToTelegram(msg, "请在命令后追加要设置的新昵称", false)
	}
	arg := args[0]
	err := catchobj.UpdateNickName(vars.DBInstance, common.UserRel{ChatId: msg.Chat.ID, UserId: msg.From.ID}, arg)
	if err == catchobj.ErrCatchObjNotFound {
		err = createCatchObj(msg, arg, "")
	}
	if err != nil {
		return err
	}
	return utils.ReplyTextToTelegram(msg, "成功设置昵称为"+arg, false)
}

// 设置emoji
func SetMyEmojiHandler(msg *models.Message) error {
	// 只在群聊中生效
	if !utils.IsGroup(msg) {
		return nil
	}

	args := utils.ParseCommandArguments(msg.Text)
	if len(args) == 0 {
		return utils.ReplyTextToTelegram(msg, "请在命令后追加要设置的新Emoji", false)
	}
	arg := args[0]
	if !isSingleEmoji(arg) {
		return utils.ReplyTextToTelegram(msg, "只能用一个emoji哦", false)
	}
	err := catchobj.UpdateEmoji(vars.DBInstance, common.UserRel{ChatId: msg.Chat.ID, UserId: msg.From.ID}, arg)
	if err == catchobj.ErrCatchObjNotFound {
		err = createCatchObj(msg, "", arg)
	}
	if err != nil {
		return err
	}
	return utils.ReplyTextToTelegram(msg, "成功设置emoji为"+arg, false)
}

// 一键同步捕捉参数
func CatchMeHereHandler(msg *models.Message) error {
	// 只在群聊中生效
	if !utils.IsGroup(msg) {
		return nil
	}
	err := catchobj.SyncLastObj(vars.DBInstance, common.UserRel{ChatId: msg.Chat.ID, UserId: msg.From.ID})
	if err != nil {
		return err
	}
	return utils.ReplyTextToTelegram(msg, "成功同步捕捉参数", false)
}

// 静默创建可抓账号
func createCatchObj(msg *models.Message, nickName, emoji string) error {
	if nickName == "" {
		nickName = tg.GetTgUserName(msg.From)
	}
	if emoji == "" {
		emoji = CATCH_DEFAULT_EMOJI
	}
	return catchobj.CreateCatchObj(vars.DBInstance, common.UserRel{ChatId: msg.Chat.ID, UserId: msg.From.ID}, nickName, emoji)
}

// 判断一个rune是否为emoji
func isEmoji(r rune) bool {
	switch {
	case r >= 0x1F600 && r <= 0x1F64F, // Emoticons
		r >= 0x1F300 && r <= 0x1F5FF, // Miscellaneous Symbols and Pictographs
		r >= 0x1F680 && r <= 0x1F6FF, // Transport and Map Symbols
		r >= 0x2600 && r <= 0x26FF,   // Misc symbols
		r >= 0x2700 && r <= 0x27BF,   // Dingbats
		r >= 0xFE00 && r <= 0xFE0F,   // Variation Selectors
		r >= 0x1F900 && r <= 0x1F9FF, // Supplemental Symbols and Pictographs
		r >= 0x1F1E6 && r <= 0x1F1FF: // Flags (iOS)
		return true
	default:
		return false
	}
}

// 判断字符串是否只有一个emoji
func isSingleEmoji(s string) bool {
	normalized := norm.NFC.String(s) // Normalize the string to ensure consistent comparison
	count := 0
	for len(normalized) > 0 {
		r, size := utf8.DecodeRuneInString(normalized)
		if !isEmoji(r) {
			return false
		}
		count++
		normalized = normalized[size:]
	}
	return count == 1
}
