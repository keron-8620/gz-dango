package auth

import (
	"context"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

// BlacklistManager 定义黑名单管理接口
type BlacklistManager interface {
	// AddToBlacklist 将token加入黑名单，duration为过期时间
	AddToBlacklist(token string, duration time.Duration) error
	// IsBlacklisted 检查token是否在黑名单中
	IsBlacklisted(token string) (bool, error)
}

// MemoryBlacklist 内存黑名单实现
type MemoryBlacklist struct {
	tokens map[string]time.Time
	mutex  sync.RWMutex
}

func NewMemoryBlacklist() *MemoryBlacklist {
	return &MemoryBlacklist{
		tokens: make(map[string]time.Time),
	}
}

func (m *MemoryBlacklist) AddToBlacklist(token string, duration time.Duration) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.tokens[token] = time.Now().Add(duration)
	return nil
}

func (m *MemoryBlacklist) IsBlacklisted(token string) (bool, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	if expireTime, exists := m.tokens[token]; exists {
		if time.Now().After(expireTime) {
			// Token已过期，从黑名单中移除
			delete(m.tokens, token)
			return false, nil
		}
		return true, nil
	}
	return false, nil
}

// RedisBlacklist Redis黑名单实现
type RedisBlacklist struct {
	client redis.Cmdable
}

func NewRedisBlacklist(client redis.Cmdable) *RedisBlacklist {
	return &RedisBlacklist{
		client: client,
	}
}

func (r *RedisBlacklist) AddToBlacklist(token string, duration time.Duration) error {
	return r.client.Set(context.Background(), "blacklist:"+token, "1", duration).Err()
}

func (r *RedisBlacklist) IsBlacklisted(token string) (bool, error) {
	result, err := r.client.Exists(context.Background(), "blacklist:"+token).Result()
	if err != nil {
		return false, err
	}
	return result > 0, nil
}
