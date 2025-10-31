package auth

import (
	"context"
	"sync"
	"time"

	"github.com/zeromicro/go-zero/core/stores/redis"
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

// RedisBlacklist 使用Redis作为后端存储实现BlacklistManager接口
// 提供分布式黑名单管理，支持超时和上下文处理
// 适用于需要共享黑名单状态的多实例应用
type RedisBlacklist struct {
	client  *redis.Redis  // Redis客户端实例
	prefix  string        // 黑名单条目的键前缀
	timeout time.Duration // 操作超时时间
}

// NewRedisBlacklist 创建一个新的RedisBlacklist实例
// client: Redis客户端实例
// prefix: 令牌的可选键前缀（如果为空则默认使用DefaultPrefix）
// timeout: 操作超时时间（如果<=0则默认使用DefaultRedisTimeout）
func NewRedisBlacklist(client *redis.Redis, prefix string, timeout time.Duration) *RedisBlacklist {
	if prefix == "" {
		prefix = DefaultPrefix
	}
	if timeout <= 0 {
		timeout = time.Duration(5) * time.Second
	}
	return &RedisBlacklist{
		client:  client,
		prefix:  prefix,
		timeout: timeout,
	}
}

// Add 将令牌添加到Redis黑名单中，duration为过期时间
// 使用Redis SET命令设置过期时间（类似SETEX行为）
// 实现上下文取消和超时处理
func (r *RedisBlacklist) Add(ctx context.Context, token string, seconds int) error {
	// 检查上下文是否已被取消
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	key := r.prefix + token
	// 为操作应用超时
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	return r.client.SetexCtx(ctx, key, "1", seconds)
}

// Remove 从Redis黑名单中移除令牌
// 使用Redis DEL命令删除键
// 此操作是幂等的 - 删除不存在的令牌不会导致错误
// 实现上下文取消和超时处理
func (r *RedisBlacklist) Remove(ctx context.Context, token string) error {
	// 检查上下文是否已被取消
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	key := r.prefix + token
	// 为操作应用超时
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	_, err := r.client.DelCtx(ctx, key)
	return err
}

// Contains 检查令牌是否存在于Redis黑名单中
// 使用Redis EXISTS命令检查键是否存在
// 实现上下文取消和超时处理
func (r *RedisBlacklist) Contains(ctx context.Context, token string) (bool, error) {
	// 检查上下文是否已被取消
	select {
	case <-ctx.Done():
		return false, ctx.Err()
	default:
	}

	key := r.prefix + token
	// 为操作应用超时
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	return r.client.ExistsCtx(ctx, key)
}
