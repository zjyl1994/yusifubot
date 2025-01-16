package startup

import (
	"errors"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/glebarez/sqlite"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	_ "github.com/joho/godotenv/autoload"
	gorm_logrus "github.com/onrik/gorm-logrus"
	"github.com/sirupsen/logrus"
	"github.com/zjyl1994/yusifubot/infra/utils"
	"github.com/zjyl1994/yusifubot/infra/vars"
	"github.com/zjyl1994/yusifubot/server/bot"
	"github.com/zjyl1994/yusifubot/server/http"
	"github.com/zjyl1994/yusifubot/service/catchgame/catchobj"
	"github.com/zjyl1994/yusifubot/service/catchgame/catchret"
	"github.com/zjyl1994/yusifubot/service/catchgame/stamina"
	"github.com/zjyl1994/yusifubot/service/tg"
	"gorm.io/gorm"
)

func Start() (err error) {
	// 加载环境变量
	vars.DebugMode, _ = strconv.ParseBool(os.Getenv("YUSIFUBOT_DEBUG"))
	if vars.DebugMode {
		logrus.SetLevel(logrus.DebugLevel)
	}

	vars.ListenAddr = utils.COALESCE(os.Getenv("YUSIFUBOT_LISTEN"), "127.0.0.1:15033")
	vars.DatabasePath = utils.COALESCE(os.Getenv("YUSIFUBOT_DATABASE_PATH"), "./yusifubot.db")

	vars.BotToken = os.Getenv("YUSIFUBOT_BOT_TOKEN")
	if vars.BotToken == "" {
		return errors.New("YUSIFUBOT_BOT_TOKEN is not set")
	}
	// 初始化数据库
	vars.DBInstance, err = gorm.Open(sqlite.Open(vars.DatabasePath), &gorm.Config{
		Logger: gorm_logrus.New(),
	})
	if err != nil {
		return err
	}
	// 初始化数据库 WAL 模式
	sqlDB, err := vars.DBInstance.DB()
	if err != nil {
		return err
	}
	_, err = sqlDB.Exec("PRAGMA journal_mode=WAL;")
	if err != nil {
		return err
	}
	err = vars.DBInstance.AutoMigrate(&tg.Chat{}, &tg.User{}, &stamina.Stamina{},
		&catchobj.CatchObj{}, &catchret.CatchRet{}, &catchret.CatchDetail{})
	if err != nil {
		return err
	}
	// 启动bot实例
	vars.BotInstance, err = tgbotapi.NewBotAPI(vars.BotToken)
	if err != nil {
		return err
	}
	vars.BotInstance.Debug = vars.DebugMode
	// 启动 bot
	go bot.Start()
	// 启动 http
	go http.Start()
	// 响应 ctrl+c
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	<-signalChan
	logrus.Infoln("Received interrupt, shutting down...")

	if err = sqlDB.Close(); err != nil {
		return err
	}
	logrus.Infoln("Shutdown complete")
	return nil
}
