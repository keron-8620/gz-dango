package auth

import (
	"context"
	"sync"
	"time"
)

// BlacklistManager 定义令牌黑名单管理接口
// 提供添加、删除和检查令牌黑名单的方法
type BlacklistManager interface {
	// Add 将令牌添加到黑名单中，duration为过期时间
	// 令牌将在指定时间后自动过期
	// 如果操作失败则返回错误
	Add(ctx context.Context, token string, seconds int) error

	// Remove 从黑名单中移除令牌
	// 如果操作失败则返回错误
	Remove(ctx context.Context, token string) error

	// Contains 检查令牌是否存在于黑名单中
	// 如果令牌在黑名单中返回true，否则返回false
	// 如果操作失败则返回错误
	Contains(ctx context.Context, token string) (bool, error)
}

const (
	// DefaultPrefix 是Redis黑名单条目使用的默认键前缀
	DefaultPrefix = "auth:blacklist:"
)

// MemoryBlacklist 使用内存存储实现BlacklistManager接口
// 适用于单实例应用或测试环境
// 通过读写锁保护实现线程安全
type MemoryBlacklist struct {
	tokens map[string]time.Time // 令牌到过期时间的映射
	mutex  sync.RWMutex         // 用于并发访问的读写锁
	prefix string               // 键前缀（主要用于与Redis实现保持一致）
}

// NewMemoryBlacklist 创建一个新的MemoryBlacklist实例
// prefix: 令牌的可选键前缀（可以为空）
func NewMemoryBlacklist(prefix string) *MemoryBlacklist {
	return &MemoryBlacklist{
		tokens: make(map[string]time.Time),
		prefix: prefix,
	}
}

// Add 将令牌添加到黑名单中，duration为过期时间
// 令牌在指定时间后将被视为过期
// 实现上下文取消和超时处理
func (m *MemoryBlacklist) Add(ctx context.Context, token string, seconds int) error {
	// 检查上下文是否已被取消
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.tokens[token] = time.Now().Add(time.Duration(seconds) * time.Second)
	return nil
}

// Remove 从黑名单中移除令牌
// 此操作是幂等的 - 删除不存在的令牌不会导致错误
// 实现上下文取消处理
func (m *MemoryBlacklist) Remove(ctx context.Context, token string) error {
	// 检查上下文是否已被取消
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	delete(m.tokens, token)
	return nil
}

// Contains 检查令牌是否存在于黑名单中且未过期
// 如果令牌存在但已过期，将自动从黑名单中移除
// 实现上下文取消处理
func (m *MemoryBlacklist) Contains(ctx context.Context, token string) (bool, error) {
	// 检查上下文是否已被取消
	select {
	case <-ctx.Done():
		return false, ctx.Err()
	default:
	}

	m.mutex.RLock()
	defer m.mutex.RUnlock()

	if expireTime, exists := m.tokens[token]; exists {
		// 检查令牌是否已过期
		if time.Now().After(expireTime) {
			// 令牌已过期，从黑名单中删除
			delete(m.tokens, token)
			return false, nil
		}
		return true, nil
	}
	return false, nil
}
