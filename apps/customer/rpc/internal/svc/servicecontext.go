package svc

import (
	"context"
	"fmt"
	"time"

	"gz-dango/apps/customer/rpc/internal/config"
	"gz-dango/pkg/auth"
	"gz-dango/pkg/database"
	"gz-dango/pkg/errors"

	"github.com/casbin/casbin/v2"
	stringadapter "github.com/casbin/casbin/v2/persist/string-adapter"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"go.etcd.io/etcd/clientv3"
	"gorm.io/gorm"
)

const (
	DefaultPolicyChangeKey = "/casbin/policy-change-signal"
)

type ServiceContext struct {
	Config   config.Config
	db       *gorm.DB
	redis    *redis.Redis
	etcd     *clientv3.Client
	enforcer *auth.AuthEnforcer

	Perm   *PermissionService
	Menu   *MenuService
	Button *ButtonService
	Role   *RoleService
	User   *UserService
	Recode *RecordService
}

func NewServiceContext(c config.Config) *ServiceContext {
	logx.MustSetup(c.LogConf)
	if len(c.Etcd.Hosts) == 0 {
		panic("请配置etcd")
	}
	// 验证etcd配置
	if err := c.Etcd.Validate(); err != nil {
		logx.Errorw("etcd配置无效", logx.Field(errors.ErrKey, err))
		panic(err)
	}
	etcdCfg := clientv3.Config{
		Endpoints:            c.Etcd.Hosts,
		DialTimeout:          5 * time.Second,
		DialKeepAliveTime:    30 * time.Second,
		DialKeepAliveTimeout: 10 * time.Second,
		AutoSyncInterval:     1 * time.Minute,
		PermitWithoutStream:  true,
	}

	if c.Etcd.User != "" && c.Etcd.Pass != "" {
		etcdCfg.Username = c.Etcd.User
		etcdCfg.Password = c.Etcd.Pass
	}

	etcd, err := clientv3.New(etcdCfg)
	if err != nil {
		logx.Errorw("创建etcd客户端失败", logx.Field(errors.ErrKey, err))
	}

	redisClient := redis.MustNewRedis(c.CacheConf.RedisConf)

	db, err := database.NewGormDB(c.DBConf, database.NewGormConfig(nil))
	if err != nil {
		logx.Errorw("创建数据库连接失败", logx.Field(errors.ErrKey, err))
		panic(err)
	}

	adapter := stringadapter.NewAdapter(`p, admin, *, *`)
	enf, err := casbin.NewEnforcer("etc/model.conf", adapter)
	if err != nil {
		panic(err)
	}
	enforcer := auth.NewAuthEnforcer(enf, "")

	return &ServiceContext{
		Config:   c,
		db:       db,
		redis:    redisClient,
		etcd:     etcd,
		enforcer: enforcer,
		Perm:     NewPermissionService(db, enforcer),
		Menu:     NewMenuService(db, enforcer),
		Button:   NewButtonService(db, enforcer),
		Role:     NewRoleService(db, enforcer),
		User:     NewUserService(db, enforcer),
		Recode:   NewRecordService(db),
	}
}

func (s *ServiceContext) DB() *gorm.DB {
	return s.db
}

func (s *ServiceContext) Redis() *redis.Redis {
	return s.redis
}

func (s *ServiceContext) Etcd() *clientv3.Client {
	return s.etcd
}

func (s *ServiceContext) Enforce() *auth.AuthEnforcer {
	return s.enforcer
}

func (s *ServiceContext) Close() {
	// 关闭 Etcd 连接
	if err := s.etcd.Close(); err != nil {
		logx.Errorw("关闭 Etcd 连接失败", logx.Field(errors.ErrKey, err))
	} else {
		logx.Info("关闭 Etcd 连接成功")
	}
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
	timeout := s.Config.AuthConf.PolicyLoadTimeout
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

	logx.Info("所有策略已刷新")
	return nil
}

// WatchCasbinPolicies 监听etcd中的casbin策略变化信号
func (s *ServiceContext) WatchCasbinPolicies() {
	// 获取策略变更key，如果未配置则使用默认值
	policyChangeKey := s.Config.AuthConf.PolicyChangeKey
	if policyChangeKey == "" {
		policyChangeKey = DefaultPolicyChangeKey
	}
	go func() {
		logx.Info("开始监听casbin策略变更信号: %s", policyChangeKey)

		// 监听key的变化
		watchChan := s.etcd.Watch(context.Background(), policyChangeKey)
		for watchResp := range watchChan {
			for range watchResp.Events {
				logx.Info("收到casbin策略变更信号, 开始刷新策略...")

				// 当收到策略变更信号时，从数据库重新加载策略
				if err := s.RefreshPolicies(); err != nil {
					logx.Errorw("刷新casbin策略失败", logx.Field(errors.ErrKey, err))
				} else {
					logx.Info("casbin策略刷新成功")
				}
			}
		}
	}()
}

// NotifyPolicyChange 发送策略变更通知到etcd
func (s *ServiceContext) NotifyPolicyChange() error {
	// 获取策略变更key，如果未配置则使用默认值
	policyChangeKey := s.Config.AuthConf.PolicyChangeKey
	if policyChangeKey == "" {
		policyChangeKey = DefaultPolicyChangeKey
	}
	// 发送当前时间戳作为信号值
	value := fmt.Sprintf("%d", time.Now().UnixNano())

	// 将信号写入etcd
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := s.etcd.Put(ctx, policyChangeKey, value)
	if err != nil {
		logx.Errorw("发送策略变更信号失败", logx.Field(errors.ErrKey, err))
		return err
	}

	logx.Info("策略变更信号发送成功")
	return nil
}
