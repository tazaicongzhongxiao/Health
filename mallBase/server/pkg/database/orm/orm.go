package orm

import (
	"MyTestMall/mallBase/basics/pkg/app"
	"MyTestMall/mallBase/basics/pkg/config"
	"MyTestMall/mallBase/basics/pkg/log"
	"errors"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/dbresolver"
	"time"
)

type (
	connInfo struct {
		Datas    []string `mapstructure:"datas"`    // 指定表名
		Sources  []string `mapstructure:"sources"`  // 使用的主DNS
		Replicas []string `mapstructure:"replicas"` // 使用的从DNS
		MaxIdle  int      `mapstructure:"max_idle"` // 设置空闲连接池中连接的最大数量
		MaxOpen  int      `mapstructure:"max_open"` // 设置打开数据库连接的最大数量
	}
	configs struct {
		Default struct {
			Dns     string `mapstructure:"dns"`      // 默认DNS
			MaxIdle int    `mapstructure:"max_idle"` // 设置空闲连接池中连接的最大数量
			MaxOpen int    `mapstructure:"max_open"` // 设置打开数据库连接的最大数量
		} `mapstructure:"default"`
		Resolver []connInfo `mapstructure:"resolver"`
	}
)

var (
	orm  *gorm.DB
	conf configs
)

// Start 启动数据库
func Start() {
	err := config.Config().Bind(app.CfgName, "mysql", &conf, func() {
		conf = configs{}
		if config.Config().Bind(app.CfgName, "mysql", &conf, nil) == nil {
			_ = initStart()
		}
	})
	if err != nil {
		log.Warn("mysql 配置错误错误", err.Error())
		panic(err)
		return
	}
	err = initStart()
	return
}

func initStart() (err error) {
	orm, err = Connection(conf.Default.Dns, conf.Default.MaxIdle, conf.Default.MaxOpen)
	if err != nil {
		log.Warn(fmt.Sprintf("数据库[%s]连接错误", conf.Default.Dns), err.Error())
		panic(err)
		return
	}
	if len(conf.Resolver) > 0 {
		var resolver *dbresolver.DBResolver
		for _, info := range conf.Resolver {
			var c dbresolver.Config
			for _, v := range info.Sources {
				c.Sources = append(c.Sources, mysql.Open(v))
			}
			for _, v := range info.Replicas {
				c.Replicas = append(c.Replicas, mysql.Open(v))
			}
			if resolver == nil {
				if len(info.Datas) > 0 {
					resolver = dbresolver.Register(c, info.Datas)
				} else {
					resolver = dbresolver.Register(c)
				}
			} else {
				if len(info.Datas) > 0 {
					resolver.Register(c, info.Datas)
				} else {
					resolver.Register(c)
				}
			}
			if info.MaxIdle > 0 {
				resolver.SetMaxIdleConns(info.MaxIdle)
			}
			if info.MaxOpen > 0 {
				resolver.SetMaxOpenConns(info.MaxOpen)
			}
		}
		if err = orm.Use(resolver); err != nil {
			panic(fmt.Sprintf("failed to use plugin, got error: %v", err))
			return
		}
	}
	return err
}

func Connection(dns string, maxIdle, maxOpen int) (orm *gorm.DB, err error) {
	orm, err = gorm.Open(mysql.New(mysql.Config{
		DSN:                       dns,  // DSN data source name
		DefaultStringSize:         256,  // string 类型字段的默认长度
		DisableDatetimePrecision:  true, // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true, // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true, // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: true, // 根据当前 MySQL 版本自动配置
	}), &gorm.Config{
		SkipDefaultTransaction: true, // 禁用默认事务
	})
	if err != nil {
		log.Warn("open mysql error", err.Error())
		return orm, err
	}
	if app.Cfg.LogLevel == "DEBUG" {
		orm.Logger.LogMode(logger.Silent)
	} else if app.Cfg.LogLevel == "INFO" {
		orm.Logger.LogMode(logger.Info)
	} else if app.Cfg.LogLevel == "WARN" {
		orm.Logger.LogMode(logger.Warn)
	} else if app.Cfg.LogLevel == "ERROR" {
		orm.Logger.LogMode(logger.Error)
	}
	sqlDB, err := orm.DB()
	if err != nil {
		log.Warn("mysql connect error", err.Error())
		return orm, err
	}
	if err = sqlDB.Ping(); err != nil {
		log.Warn("mysql ping error: %s", err.Error())
		return orm, err
	}
	// 连接池里面的连接最大空闲时长
	sqlDB.SetConnMaxIdleTime(0)
	// SetMaxIdleConns 表示设置最大的可空闲连接数，该函数的作用就是保持等待连接操作状态的连接数，这个主要就是避免操作过程中频繁的获取连接，释放连接。默认情况下会保持的连接数量为2.就是说会有两个连接一直保持，不释放，等待需要使用的用户使用。
	sqlDB.SetMaxIdleConns(maxIdle)
	// SetMaxOpenConns 表示最大的连接数，这个我们不设置默认就是不限制，可以无限创建连接，问题就在数据库本身有瓶颈，无限创建，会损耗性能。所以我们要根据我们自己的数据库瓶颈情况来进行相关的设置。当出现连接数超出了我们设定的数量时候，后面的用户等待超时时间之前，有连接释放就会自动获得操作的权限，否则返回连接超时。（每个公司的使用情况不同，所以根据情况自己设定，个人建议不要采用默认无限制创建连接）
	sqlDB.SetMaxOpenConns(maxOpen)
	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	sqlDB.SetConnMaxLifetime(time.Hour)
	return orm, err
}

// Client
// @Description: 获取MYSQL
// @return *gorm.DB
func Client() *gorm.DB {
	return orm
}

// Ping
// @Description: 检查MYSQL状态失败
// @return error
func Ping() error {
	sqlDB, err := orm.DB()
	if err != nil {
		log.Warn("sql connect db server failed.", err.Error())
	}
	if err = sqlDB.Ping(); err != nil {
		log.Warn("sql ping failed.", err.Error())
	}
	return err
}

// Read
// @Description: 使用 读模式
// @return *gorm.DB
func Read() *gorm.DB {
	return Client().Clauses(dbresolver.Read)
}

// Write
// @Description: 使用 写模式
// @return *gorm.DB
func Write() *gorm.DB {
	return Client().Clauses(dbresolver.Write)
}

// IsNotFound
// @Description: 检查 ErrRecordNotFound 错误
// @param err
// @return bool
func IsNotFound(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}

// Collection
// @Description: Collection 得到一个mysql操作对象
// @return res
func Collection() (res *CollectionInfo) {
	return &CollectionInfo{client: orm}
}
