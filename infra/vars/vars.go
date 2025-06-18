package vars

import (
	"math/rand/v2"

	tgbot "github.com/go-telegram/bot"
	"gorm.io/gorm"
)

var (
	DebugMode  bool
	ListenAddr string

	BotToken    string
	BotInstance *tgbot.Bot

	DatabasePath string
	DBInstance   *gorm.DB

	AdminUserId string

	RNG *rand.ChaCha8
)
