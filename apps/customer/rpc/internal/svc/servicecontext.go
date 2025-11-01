package svc

import (
	"context"
	"time"

	"gz-dango/apps/customer/rpc/internal/config"
	"gz-dango/apps/customer/rpc/internal/models"
	"gz-dango/pkg/auth"
	"gz-dango/pkg/database"
	"gz-dango/pkg/errors"

	"github.com/casbin/casbin/v2"
	stringadapter "github.com/casbin/casbin/v2/persist/string-adapter"
	"github.com/google/uuid"
	goReids "github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"gorm.io/gorm"
)

const (
	DefaultPolicyChangeKey = "/casbin/policy-change-signal"
)

type ServiceContext struct {
	Config     config.Config
	db         *gorm.DB
	redis      *redis.Redis
	goredis    *goReids.Client
	enforcer   *auth.AuthEnforcer
	instanceID string

	Perm   *PermissionService
	Menu   *MenuService
	Button *ButtonService
	Role   *RoleService
	User   *UserService
	Recode *RecordService
}

func NewServiceContext(c config.Config) *ServiceContext {
	redisClient := redis.MustNewRedis(c.Cache)

	goredisClient := goReids.NewClient(&goReids.Options{
		Addr:     c.Cache.Host,
		Password: c.Cache.Pass,
	})

	db, err := database.NewGormDB(c.Database, database.NewGormConfig(nil))
	if err != nil {
		logx.Errorw("创建数据库连接失败", logx.Field(errors.ErrKey, err))
		panic(err)
	}
	if err := db.AutoMigrate(
		&models.PermissionModel{},
		&models.MenuModel{},
		&models.ButtonModel{},
		&models.RoleModel{},
		&models.UserModel{},
		&models.LoginRecordModel{},
	); err != nil {
		logx.Errorw("数据库自动迁移失败", logx.Field(errors.ErrKey, err))
		panic(err)
	}

	adapter := stringadapter.NewAdapter(`p, admin, *, *`)
	enf, err := casbin.NewEnforcer("etc/model.conf", adapter)
	if err != nil {
		panic(err)
	}
	enforcer := auth.NewAuthEnforcer(enf, "")
	enforcer.SetBlacklist(
		auth.NewRedisBlacklist(
			redisClient,
			c.Security.JwtBlacklistPrefix,
			time.Duration(c.Security.TokenExpireMinutes)*time.Minute,
		),
	)
	return &ServiceContext{
		Config:     c,
		db:         db,
		redis:      redisClient,
		goredis:    goredisClient,
		enforcer:   enforcer,
		instanceID: uuid.New().String(),
		Perm:       NewPermissionService(db, enforcer),
		Menu:       NewMenuService(db, enforcer),
		Button:     NewButtonService(db, enforcer),
		Role:       NewRoleService(db, enforcer),
		User:       NewUserService(db, enforcer),
		Recode:     NewRecordService(db),
	}
}

func (s *ServiceContext) DB() *gorm.DB {
	return s.db
}

func (s *ServiceContext) Redis() *redis.Redis {
	return s.redis
}

func (s *ServiceContext) GoRedis() *goReids.Client {
	return s.goredis
}

func (s *ServiceContext) Enforce() *auth.AuthEnforcer {
	return s.enforcer
}

func (s *ServiceContext) Close() {
	// 关闭数据库连接
	conn, err := s.db.DB()
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
	timeout := s.Config.Security.PolicyLoadTimeout
	if timeout <= 0 {
		timeout = 30 * time.Second
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

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
	return nil
}

// WatchCasbinPolicies 监听casbin策略变化信号
func (s *ServiceContext) WatchCasbinPolicies() {
	// 获取策略变更key，如果未配置则使用默认值
	policyChangeKey := s.Config.Security.PolicyChangeKey
	if policyChangeKey == "" {
		policyChangeKey = DefaultPolicyChangeKey
	}
	go func() {
		logx.Infow("开始监听casbin策略变更信号", logx.Field("channel", policyChangeKey))

		// 使用Redis的Pub/Sub功能监听键空间事件
		pubsub := s.goredis.Subscribe(context.Background(), policyChangeKey)
		defer pubsub.Close()

		ch := pubsub.Channel()
		for msg := range ch {
			if msg.Payload == s.instanceID {
				continue
			}
			logx.Infow("收到Redis同步信号: 开始同步策略", logx.Field("message", msg.Payload))
			if err := s.RefreshPolicies(); err != nil {
				logx.Errorw("策略同步失败", logx.Field(errors.ErrKey, err))
			}
			logx.Infow("策略已同步完成")
		}
	}()
}

// NotifyPolicyChange 发送策略变更通知
func (s *ServiceContext) NotifyPolicyChange() error {
	// 获取策略变更key，如果未配置则使用默认值
	policyChangeKey := s.Config.Security.PolicyChangeKey
	if policyChangeKey == "" {
		policyChangeKey = DefaultPolicyChangeKey
	}

	// 将信号写入etcd
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 将信号发布到Redis频道
	if err := s.goredis.Publish(ctx, policyChangeKey, s.instanceID).Err(); err != nil {
		logx.Errorw("发送Redis同步信号失败", logx.Field(errors.ErrKey, err))
		return err
	}

	logx.Info("策略变更信号发送成功")
	return nil
}
