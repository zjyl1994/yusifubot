package sign

import (
	"encoding/json"
	"errors"
	"math/rand/v2"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/zjyl1994/yusifubot/infra/vars"
	"github.com/zjyl1994/yusifubot/service/catchgame/common"
	"github.com/zjyl1994/yusifubot/service/catchgame/stamina"
	"github.com/zjyl1994/yusifubot/service/configure"
	"gorm.io/gorm"
)

func userSign(user common.UserRel) error {
	signLock.Lock(user)
	defer signLock.Unlock(user)

	var m CatchSign
	err := vars.DBInstance.Where("chat_id = ? AND user_id = ?", user.ChatId, user.UserId).First(&m).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		m.ChatId = user.ChatId
		m.UserId = user.UserId
		err = vars.DBInstance.Create(&m).Error
		if err != nil {
			return err
		}
	}

	now := time.Now()
	signYear, signMonth, signDay := time.Unix(m.SignTime, 0).Date()
	todayYear, todayMonth, todayDay := now.Date()

	if signYear == todayYear && signMonth == todayMonth && signDay == todayDay {
		return ErrAlreadySigned
	}

	m.SignTime = now.Unix()

	return vars.DBInstance.Save(&m).Error
}

func Sign(user common.UserRel) (int64, error) {
	var minSp, maxSp int64
	minSp = defaultSpMin
	maxSp = defaultSpMax
	// 检查是否有sp奖励内容
	signSpCfgStr, err := configure.Get("catch_sign_sp", "[]")
	if err == nil {
		var activitys []signActivityConfig
		err = json.Unmarshal([]byte(signSpCfgStr), &activitys)
		if err == nil {
			for _, act := range activitys {
				startTime, err := time.Parse(time.DateTime, act.StartTime)
				if err != nil {
					continue
				}
				endTime, err := time.Parse(time.DateTime, act.EndTime)
				if err != nil {
					continue
				}
				now := time.Now()
				if now.After(startTime) && now.Before(endTime) {
					logrus.Debugln("SignAct", startTime, endTime, act.Min, act.Max)
					minSp = act.Min
					maxSp = act.Max
				}
			}
		}
	}
	// 计算签到得到的体力
	actualSp := rand.Int64N(maxSp-minSp) + minSp
	// 签到
	if err := userSign(user); err != nil {
		return 0, err
	}
	// 给予对应奖励
	err = stamina.AddStaminPoint(user, actualSp)
	return int64(actualSp), err
}
