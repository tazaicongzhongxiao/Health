package config

import (
	"github.com/spf13/viper"
	"mallBase/basics/pkg/app"
)

type Configuration struct {
	configs map[string]*configName
}

type configName struct {
	uptime int64        // 更新时间
	viper  *viper.Viper // 配置
	fn     []func()     // 配置变更后执行
}

var config = new(Configuration).init()

func (c *Configuration) init() *Configuration {
	c.configs = make(map[string]*configName)
	return c
}

func Config() *Configuration {
	return config
}

func (c *Configuration) getRemoteConf(fileName string) (conf *viper.Viper, err error) {
	if objVal, ok := c.configs[fileName]; ok {
		return objVal.viper, nil
	} else {
		var cache string
		if app.Cfg.BasePath == "" { //

		}
	}
}
