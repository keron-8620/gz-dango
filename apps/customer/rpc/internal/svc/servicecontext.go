package svc

import (
	"context"
	"gz-dango/apps/customer/rpc/internal/config"
	"gz-dango/pkg/auth"
	"gz-dango/pkg/database"
	"gz-dango/pkg/errors"
	"time"

	"github.com/casbin/casbin/v2"
	stringadapter "github.com/casbin/casbin/v2/persist/string-adapter"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config config.Config
	DB     *gorm.DB
	Redis  *redis.Redis

	Perm   *PermissionService
	Menu   *MenuService
	Button *ButtonService
	Role   *RoleService
	User   *UserService
	Recode *RecordService
}

func NewServiceContext(c config.Config) *ServiceContext {
	logx.MustSetup(c.LogConf)
	redisClient := redis.MustNewRedis(c.CacheConf.RedisConf)
	db, err := database.NewGormDB(c.DBConf, database.NewGormConfig(nil))
	if err != nil {
		panic(err)
	}

	adapter := stringadapter.NewAdapter(`p, admin, *, *`)
	enf, err := casbin.NewEnforcer("etc/model.conf", adapter)
	if err != nil {
		panic(err)
	}
	enforcer := auth.NewAuthEnforcer(enf, "")
	if redisClient != nil {
		bl := NewRedisBlacklist(redisClient, auth.DefaultPrefix, DefaultRedisTimeout)
		enforcer.SetBlacklist(bl)
	}

	return &ServiceContext{
		Config: c,
		DB:     db,
		Redis:  redisClient,
		Perm:   NewPermissionService(db, enforcer),
		Menu:   NewMenuService(db, enforcer),
		Button: NewButtonService(db, enforcer),
		Role:   NewRoleService(db, enforcer),
		User:   NewUserService(db, enforcer),
		Recode: NewRecordService(db),
	}
}

func (s *ServiceContext) Close() {
	// 关闭数据库连接
	conn, err := s.DB.DB()
	if err != nil {
		logx.Errorw("获取数据库连接失败", logx.Field(errors.ErrKey, err))
	}
	if err = conn.Close(); err != nil {
		logx.Errorw("关闭数据库连接失败", logx.Field(errors.ErrKey, err))
	} else {
		logx.Info("关闭数据库连接成功")
	}
}

// 提供一个 HTTP 接口用于刷新 Casbin 策略
func (s *ServiceContext) RefreshPolicies() error {
	ctx := context.Background()

	if err := s.Perm.LoadPolicies(ctx); err != nil {
		return err
	}

	if err := s.Menu.LoadPolicies(ctx); err != nil {
		return err
	}

	if err := s.Button.LoadPolicies(ctx); err != nil {
		return err
	}

	if err := s.Role.LoadPolicies(ctx); err != nil {
		return err
	}

	logx.Info("所有策略已刷新")
	return nil
}

const (
	// DefaultRedisTimeout 是Redis操作的默认超时时间
	DefaultRedisTimeout = time.Duration(5) * time.Second
)

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
		prefix = auth.DefaultPrefix
	}
	if timeout <= 0 {
		timeout = DefaultRedisTimeout
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
