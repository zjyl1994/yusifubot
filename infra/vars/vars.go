package vars

import (
	"math/rand/v2"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"
)

var (
	DebugMode  bool
	ListenAddr string

	BotToken    string
	BotInstance *tgbotapi.BotAPI

	DatabasePath string
	DBInstance   *gorm.DB

	AdminUserId string

	RNG *rand.ChaCha8
)
