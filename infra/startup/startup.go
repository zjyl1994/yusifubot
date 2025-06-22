package startup

import (
	"context"
	crand "crypto/rand"
	"errors"
	"math/rand/v2"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/glebarez/sqlite"
	tgbot "github.com/go-telegram/bot"
	"github.com/gofiber/fiber/v2"
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
	"github.com/zjyl1994/yusifubot/service/config"
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

	vars.AdminUserId, err = strconv.ParseInt(os.Getenv("YUSIFUBOT_ADMIN_USER_ID"), 10, 64)
	if err != nil {
		return err
	}
	vars.AdminToken = os.Getenv("YUSIFUBOT_ADMIN_TOKEN")
	vars.ReplicateToken = os.Getenv("YUSIFUBOT_REPLICATE_TOKEN")
	vars.ReplicateCooldown = utils.NewCooldownManager()

	// 初始化独立随机数器
	var randSeed [32]byte
	_, err = crand.Read(randSeed[:])
	if err != nil {
		return err
	}
	vars.RNG = rand.NewChaCha8(randSeed)

	// 初始化数据库
	vars.DBInstance, err = gorm.Open(sqlite.Open(vars.DatabasePath), &gorm.Config{
		Logger:         gorm_logrus.New(),
		TranslateError: true,
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
	err = vars.DBInstance.AutoMigrate(&tg.Chat{}, &tg.User{}, &stamina.Stamina{}, &catchobj.CatchObj{}, &catchret.CatchResult{}, &config.ChatConfig{})
	if err != nil {
		return err
	}
	// 初始化AdminApi
	adminApi := fiber.New()
	http.AdminApi(adminApi)
	// 启动bot实例
	botOpts := []tgbot.Option{
		tgbot.WithDefaultHandler(bot.Handler),
	}
	vars.BotInstance, err = tgbot.New(vars.BotToken, botOpts...)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// 启动 bot
	go vars.BotInstance.Start(ctx)
	// 启动 admin api
	go func() {
		if err = adminApi.Listen(vars.ListenAddr); err != nil {
			logrus.Errorln("Admin API failed:", err)
		}
	}()
	// 响应 ctrl+c
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	<-signalChan
	logrus.Infoln("Received interrupt, shutting down...")

	if err = adminApi.ShutdownWithTimeout(3 * time.Second); err != nil {
		return err
	}
	if err = sqlDB.Close(); err != nil {
		return err
	}
	logrus.Infoln("Shutdown complete")
	return nil
}
