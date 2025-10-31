package config

import (
	"time"

	"gz-dango/pkg/database"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	LogConf   logx.LogConf
	CacheConf cache.NodeConf
	DBConf    database.DBConf
	AuthConf  AuthConfig
}

type AuthConfig struct {
	JwtSecret         string
	PolicyLoadTimeout time.Duration
	PolicyChangeKey   string
}
