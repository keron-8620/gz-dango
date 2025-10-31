// Package database 提供数据库连接和配置功能
// 支持多种数据库类型包括MySQL、PostgreSQL、SQLite、SQLServer和OpenGauss
package database

import (
	"io"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"gz-dango/pkg/database/driver/opengauss"
)

// DBConf 数据库配置结构体，用于配置数据库连接参数
type DBConf struct {
	Type            string `yaml:"type" json:"type"`                       // 数据库类型，支持 mysql, postgres, sqlite, sqlserver, opengauss
	Dns             string `yaml:"dns" json:"dns"`                         // 数据库连接字符串
	MaxIdleConns    int    `yaml:"maxIdleConns" json:"maxIdleConns"`       // 最大空闲连接数
	MaxOpenConns    int    `yaml:"maxOpenConns" json:"maxOpenConns"`       // 最大打开连接数
	ConnMaxLifetime int    `yaml:"connMaxLifetime" json:"connMaxLifetime"` // 连接最大生命周期(秒)
	ConnMaxIdleTime int    `yaml:"connMaxIdleTime" json:"connMaxIdleTime"` // 连接最大空闲时间(秒)
}

// NewGormConfig 创建 gorm 配置
func NewGormConfig(w io.Writer) *gorm.Config {
	var gc gorm.Config
	gc.SkipDefaultTransaction = false
	gc.FullSaveAssociations = false
	gc.TranslateError = true
	if w != nil {
		dbLog := log.New(w, " ", log.LstdFlags)
		gc.Logger = logger.New(dbLog, logger.Config{
			SlowThreshold:             time.Second,
			Colorful:                  false,
			IgnoreRecordNotFoundError: false,
			ParameterizedQueries:      false,
			LogLevel:                  logger.Info,
		})
	}
	return &gc
}

// InitGormDB 初始化GORM数据库连接
// c: 数据库配置信息
// gc: GORM配置信息
// 返回GORM数据库实例和可能的错误
func NewGormDB(conf DBConf, gc *gorm.Config) (*gorm.DB, error) {
	var (
		db      *gorm.DB // GORM数据库实例
		openErr error    // 数据库打开错误
	)

	// 根据数据库类型选择相应的驱动并建立连接
	switch conf.Type {
	case "mysql":
		db, openErr = gorm.Open(mysql.Open(conf.Dns), gc)
	case "postgres":
		db, openErr = gorm.Open(postgres.Open(conf.Dns), gc)
	case "sqlite":
		db, openErr = gorm.Open(sqlite.Open(conf.Dns), gc)
	case "sqlserver":
		db, openErr = gorm.Open(sqlserver.Open(conf.Dns), gc)
	case "opengauss":
		db, openErr = gorm.Open(opengauss.Open(conf.Dns), gc)
	default:
		// 不支持的数据库驱动类型
		return nil, gorm.ErrUnsupportedDriver
	}

	// 如果数据库连接打开失败，返回错误
	if openErr != nil {
		return nil, openErr
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	// 设置连接池参数
	sqlDB.SetMaxIdleConns(conf.MaxIdleConns)                                    // 最大空闲连接数
	sqlDB.SetMaxOpenConns(conf.MaxOpenConns)                                    // 最大打开连接数
	sqlDB.SetConnMaxLifetime(time.Duration(conf.ConnMaxLifetime) * time.Second) // 连接最大生命周期
	sqlDB.SetConnMaxIdleTime(time.Duration(conf.ConnMaxIdleTime) * time.Second) // 空闲连接最大存活时间
	return db, nil
}

// CloseGormDB 关闭GORM数据库连接
// db: GORM数据库实例
// 返回关闭操作可能产生的错误
func CloseGormDB(db *gorm.DB) error {
	// 检查数据库实例是否有效
	if db == nil {
		return gorm.ErrInvalidDB
	}

	// 获取底层数据库连接
	sqlDB, err := db.DB()
	if err != nil {
		return gorm.ErrInvalidDB
	}

	// 关闭数据库连接
	return sqlDB.Close()
}
