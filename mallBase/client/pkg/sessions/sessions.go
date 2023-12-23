package sessions

import (
	"MyTestMall/mallBase/basics/pkg/app"
	"MyTestMall/mallBase/basics/pkg/config"
	"MyTestMall/mallBase/basics/pkg/log"
	"github.com/gin-contrib/sessions"
	redisSession "github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"strconv"
)

type configs struct {
	Key         string `mapstructure:"key"`
	Name        string `mapstructure:"name"`
	Domain      string `mapstructure:"domain"`
	Addr        string `mapstructure:"addr"`
	Password    string `mapstructure:"password"`
	Db          int    `mapstructure:"db"`
	PoolSize    int    `mapstructure:"pool_size"`
	MinIdleConn int    `mapstructure:"min_idle_conn"`
}

var conf configs

// Inject 启动session服务, 在自定义的路由代码中调用, 传入 *gin.Engine 对象
func Inject(engine *gin.Engine) {
	if err := config.Config().Bind(app.CfgName, "sessions", &conf, func() {
		if config.Config().Bind(app.CfgName, "sessions", &conf, nil) == nil {
			initStart(engine)
		}
	}); err == nil {
		initStart(engine)
	}
	return
}

func initStart(engine *gin.Engine) gin.IRoutes {
	store, err := redisSession.NewStoreWithDB(conf.PoolSize, "tcp", conf.Addr, conf.Password, strconv.Itoa(conf.Db), []byte(conf.Key))
	if err != nil {
		log.Error("sessions new", err.Error())
		return engine
	}
	store.Options(sessions.Options{MaxAge: 3600, Path: "/", Domain: conf.Domain, HttpOnly: true})
	return engine.Use(sessions.Sessions(conf.Name, store))
}

// Get 获取指定session
func Get(c *gin.Context, key string) string {
	sess := sessions.Default(c)
	val := sess.Get(key)
	if val != nil {
		return val.(string)
	}
	return ""
}

// Set 设置session
func Set(c *gin.Context, key, val string) {
	sess := sessions.Default(c)
	sess.Set(key, val)
	_ = sess.Save()
}

// Del 删除指定session
func Del(c *gin.Context, key string) {
	sess := sessions.Default(c)
	sess.Delete(key)
	_ = sess.Save()
}
