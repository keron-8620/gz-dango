package config

import (
	"time"

	"gz-dango/pkg/database"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	ServiceLog logx.LogConf
	Database   database.DBConf
	Cache      redis.RedisConf
	Security   SecurityConfig
}

type SecurityConfig struct {
	JwtSecret          string
	PolicyLoadTimeout  time.Duration
	PolicyChangeKey    string
	JwtBlacklistPrefix string
	CheckTimestamp     bool
	TimestampRange     int
	TokenExpireMinutes int
	LoginFailMaxTimes  int
	PasswordStrength   int
}
