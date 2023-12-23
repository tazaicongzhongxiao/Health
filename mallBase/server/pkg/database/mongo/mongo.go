package mongo

import (
	"MyTestMall/mallBase/basics/pkg/app"
	"MyTestMall/mallBase/basics/pkg/config"
	"MyTestMall/mallBase/basics/pkg/log"
	mongodb "MyTestMall/mallBase/basics/pkg/mongo"
	"MyTestMall/mallBase/server/pkg/database"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

var (
	client *mongo.Client
	conf   configs
)

type (
	configs struct {
		URL             string `mapstructure:"url"`
		Database        string `mapstructure:"database"`
		MaxConnIdleTime int    `mapstructure:"max_conn_idle_time"`
		MaxPoolSize     int    `mapstructure:"max_pool_size"`
		ReplSetName     string `mapstructure:"repl_set_name"`
		Direct          bool   `mapstructure:"direct"` // 强制此客户端
		AuthSource      string `mapstructure:"auth_source"`
		Username        string `mapstructure:"username"`
		Password        string `mapstructure:"password"`
	}
)

// Start 启动 mongo
func Start(index map[string]interface{}) {
	err := config.Config().Bind(app.CfgName, "mongo", &conf, func() {
		conf = configs{}
		if config.Config().Bind(app.CfgName, "mongo", &conf, nil) == nil {
			if db, err := initStart(); err == nil {
				client = db
			}
		}
	})
	if err != nil {
		log.Error("mongo 配置错误", err.Error())
		panic(err)
		return
	}
	if client, err = initStart(); err != nil {
		log.Error("mongo initStart", err.Error())
		panic(err)
		return
	}
	// 开始创建索引
	if len(index) > 0 {
		log.Info(fmt.Sprintf("mongo 开始创建索引: %d 条", len(index)), nil)
		opts := options.CreateIndexes().SetMaxTime(3600 * time.Second)
		db := client.Database(conf.Database)
		for name1, val1 := range index {
			go func(name string, val interface{}) {
				var data []mongodb.IndexData
				if err = app.Unmarshal(val, &data); err == nil {
					if _, err = db.Collection(name).Indexes().CreateMany(context.Background(), GetIndexModel(data), opts); err != nil {
						log.Info(fmt.Sprintf("mongo 创建索引 %s", name), err.Error())
					}
				}
			}(name1, val1)
		}
	}
}

func initStart() (db *mongo.Client, err error) {
	return Connection(conf)
}

func Client() *mongo.Client {
	return client
}

// Database 获取数据库连接
func Database(name ...string) *CollectionInfo {
	data := conf.Database
	if len(name) == 1 {
		data = name[0]
	}
	return GetDataBase(client, data)
}

// Collection 得到一个mongo操作对象
func Collection(table database.Table) *CollectionInfo {
	return GetCollection(client, conf.Database, table)
}
