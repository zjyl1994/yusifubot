package utils

import (
	"sync"
	"time"
)

// CooldownManager 管理器
type CooldownManager struct {
	mu        sync.RWMutex
	cooldowns map[string]time.Time // key -> cooldownEndTime
}

// NewCooldownManager 初始化
func NewCooldownManager() *CooldownManager {
	return &CooldownManager{
		cooldowns: make(map[string]time.Time),
	}
}

// SetCooldown 设置指定 key 的冷却时间
func (cm *CooldownManager) SetCooldown(key string, duration time.Duration) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.cooldowns[key] = time.Now().Add(duration)
}

// IsReady 检查 key 是否已经冷却好
func (cm *CooldownManager) IsReady(key string) bool {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	end, exists := cm.cooldowns[key]
	return !exists || time.Now().After(end)
}
// CheckAndSetCooldown 检查是否冷却完成，如果完成则设置新的冷却时间并返回true
func (cm *CooldownManager) CheckAndSetCooldown(key string, duration time.Duration) bool {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	
	end, exists := cm.cooldowns[key]
	if !exists || time.Now().After(end) {
		cm.cooldowns[key] = time.Now().Add(duration)
		return true
	}
	return false
}

// RemainingTime 返回剩余冷却时间，如果已经冷却完成返回 0
func (cm *CooldownManager) RemainingTime(key string) time.Duration {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	end, exists := cm.cooldowns[key]
	if !exists {
		return 0
	}
	remaining := time.Until(end)
	if remaining < 0 {
		return 0
	}
	return remaining
}

// RemoveCooldown 移除某个冷却项（可选）
func (cm *CooldownManager) RemoveCooldown(key string) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	delete(cm.cooldowns, key)
}

// Optional: 清理已过期的冷却（节省内存）
func (cm *CooldownManager) CleanupExpired() {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	now := time.Now()
	for key, end := range cm.cooldowns {
		if now.After(end) {
			delete(cm.cooldowns, key)
		}
	}
}
