package dredis

import (
	"MyTestMall/mallBase/basics/pkg/app"
	"MyTestMall/mallBase/basics/pkg/config"
	"MyTestMall/mallBase/basics/pkg/log"
	"crypto/tls"
	"github.com/go-redis/redis/v8"
	"time"
)

type Configs struct {
	Addr         []string `mapstructure:"addr"`          // 集群节点地址
	Password     string   `mapstructure:"password"`      // REDIS 密码
	Db           int      `mapstructure:"db"`            // 非集群指定 DBID
	PoolSize     int      `mapstructure:"pool_size"`     // 连接池大小 默认每CPU10个连接
	MinIdleConn  int      `mapstructure:"min_idle_conn"` // 非集群指定 最小空闲连接数
	MaxRedirects int      `mapstructure:"max_redirects"` // 集群指定 当遇到网络错误或者MOVED/ASK重定向命令时，最多重试几次，默认8
	ReadOnly     bool     `mapstructure:"read_only"`     // 集群指定 置为true则允许在从节点上执行只含读操作的命令
	Cluster      bool     `mapstructure:"cluster"`       // 是否为集群链接模式
	TLSConfig    *tls.Config
}

var (
	// client redis连接资源
	client    Driver
	RedisConf Configs
)

type Driver interface {
	Ping() error
	Close() error
	CacheGet(key string, expiration time.Duration, f func() (res string, err error)) (string, error)
	Client() (bool, *redis.Client, *redis.ClusterClient)
	Get(key string) (string, error)
	Set(key string, value interface{}, expiration time.Duration) (err error)
	Del(key string) (int64, error)
	SetNX(key string, value interface{}, expiration time.Duration) (bool, error)
	LPush(key string, values ...interface{}) (int64, error)
	LPushX(key string, values ...interface{}) (int64, error)
	Eval(script string, keys []string, args ...interface{}) (interface{}, error)
	Do(args ...interface{}) (interface{}, error)
	Expire(key string, expiration time.Duration) (bool, error)
	ExpireAt(key string, ttl time.Time) (bool, error)
	TTL(key string) (time.Duration, error)
	Exists(keys ...string) bool
	Incr(key string) int64
	ZAdd(key string, score float64, member string) (reply int64, err error)
	ZIncrby(key, member string, increment float64) (float64, error)
	ZRange(key string, start, end int64) ([]redis.Z, error)
	ZRem(key string, member string) (reply int64, err error)
	ZCard(key string) int64
	ZScore(key string, member string) (float64, error)
	BRPop(timeout time.Duration, keys ...string) ([]string, error)
	HMGet(key string, fields ...string) ([]interface{}, error)
	HGetAll(key string) (map[string]string, error)
	HIncrBy(key, field string, incr int64) (int64, error)
	HKeys(key string) ([]string, error)
	HDel(key string, fields ...string) (int64, error)
	ZRangeByScore(key string, opt *redis.ZRangeBy) ([]string, error)
	ZRemRangeByScore(key, min, max string) (int64, error)
	ZCount(key, min, max string) (int64, error)
	HSet(key string, values ...interface{}) (int64, error)
	Scan(cursor uint64, match string, count int64) ([]string, uint64, error)
	HScan(key string, cursor uint64, match string, count int64) ([]string, uint64, error)
	SAdd(key string, args ...interface{}) (int64, error)
	SCard(key string) (int64, error)
	SIsMember(key string, member interface{}) bool
	SRandMember(key string) string
	SMembers(key string) ([]string, error)
	SRem(key string, args ...interface{}) (int64, error)
	Subscribe(channels ...string) *redis.PubSub
	PSubscribe(channels ...string) *redis.PubSub
	Publish(channel string, message string) (int64, error)
}

// Start 启动redis
func Start() {
	if err := config.Config().Bind(app.CfgName, "redis", &RedisConf, func() {
		if config.Config().Bind(app.CfgName, "redis", &RedisConf, nil) == nil {
			initialize(RedisConf)
		}
	}); err != nil {
		log.Error("redis 配置错误: %s", err.Error())
		panic(err)
		return
	}
	initialize(RedisConf)
	return
}

func Get() (db Driver) {
	return client
}

func initialize(conf Configs) {
	var err error
	client, err = Connection(conf)
	if err != nil {
		log.Warn("Redis 启动错误", err.Error())
		panic(err)
	} else {
		err = client.Ping()
		if err != nil {
			panic(err)
		}
		log.Info("Redis 连接完成", conf.Addr)
	}
	return
}

func Connection(c Configs) (db Driver, err error) {
	if c.Cluster {
		opts := &redis.ClusterOptions{
			Addrs:    c.Addr,
			Password: c.Password,
			ReadOnly: c.ReadOnly,
		}
		if c.MaxRedirects > 0 {
			opts.MaxRedirects = c.MaxRedirects
		}
		if c.TLSConfig != nil {
			opts.TLSConfig = c.TLSConfig
		}
		if c.PoolSize > 0 {
			opts.PoolSize = c.PoolSize
		}
		if c.MinIdleConn > 0 {
			opts.MinIdleConns = c.MinIdleConn
		}
		db, err = NewClusterClient(opts)
	} else {
		opts := &redis.Options{
			Addr:     c.Addr[0],
			Password: c.Password,
			DB:       c.Db,
		}
		if c.TLSConfig != nil {
			opts.TLSConfig = c.TLSConfig
		}
		if c.PoolSize > 0 {
			opts.PoolSize = c.PoolSize
		}
		if c.MinIdleConn > 0 {
			opts.MinIdleConns = c.MinIdleConn
		}
		db, err = NewClient(opts)
	}
	return db, err
}
