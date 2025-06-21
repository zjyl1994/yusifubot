package vars

import (
	"time"
)

type CooldownManager interface {
	CheckAndSetCooldown(key string, duration time.Duration) bool
	CleanupExpired()
}
