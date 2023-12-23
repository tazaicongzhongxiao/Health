package config

import (
	"MyTestMall/mallBase/basics/pkg/app"
	"fmt"
	"github.com/TarsCloud/TarsGo/tars"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"time"
)

// Configuration 应用配置
type Configuration struct {
	configs map[string]*configName
}

type configName struct {
	uptime int64        // 配置最后更新时间
	viper  *viper.Viper // 配置数据
	fn     []func()     // 配置变更后执行函数
}

var (
	config = new(Configuration).init()
)

// Config 得到config对象
func Config() *Configuration {
	return config
}

func (c *Configuration) init() *Configuration {
	c.configs = make(map[string]*configName)
	return c
}

// 获取 tars 远程配置
func (c *Configuration) getRemoteConf(filename string) (conf *viper.Viper, err error) {
	if objVal, ok := c.configs[filename]; !ok {
		var cache string
		if app.Cfg.Config == "" {
			// 读取本地配置
			if bytes, err := ioutil.ReadFile(filename); err != nil {
				return conf, err
			} else {
				cache = string(bytes)
			}
			if err = c.configSetup(filename, ".", cache); err != nil {
				return conf, err
			}
		} else {
			// 读取远程配置
			remoteConf := tars.NewRConf(app.Cfg.App, app.Cfg.Server, app.Cfg.BasePath)
			cache, err = remoteConf.GetConfig(filename)
			if cache == "" || err != nil {
				return conf, fmt.Errorf("远程配置错误")
			}
			if err = c.configSetup(filename, app.Cfg.BasePath, cache); err != nil {
				return conf, err
			}
		}
		return c.configs[filename].viper, nil
	} else {
		return objVal.viper, nil
	}
}

func (c *Configuration) configSetup(filename, dir, config string) (err error) {
	c.configs[filename] = &configName{
		viper: viper.New(),
	}
	c.configs[filename].viper.SetConfigFile(dir + "/" + filename)
	c.configs[filename].viper.SetConfigType(strings.ToLower(strings.Trim(path.Ext(filename), ".")))
	if err = c.configs[filename].viper.ReadConfig(strings.NewReader(os.ExpandEnv(config))); err != nil {
		return err
	}
	c.configs[filename].viper.WatchConfig()
	c.configs[filename].viper.OnConfigChange(func(e fsnotify.Event) {
		if c.configs[filename].uptime < time.Now().Unix() {
			c.configs[filename].uptime = time.Now().Unix()
			fmt.Println("监听到文件改变！", e.Name)
			for _, val := range c.configs[filename].fn {
				go val()
			}
		}
	})
	return err
}

func (c *Configuration) Bind(filename string, key string, obj interface{}, run func()) (err error) {
	objNode, err := c.getRemoteConf(filename)
	if err == nil {
		if key != "" {
			if objVal := objNode.Sub(key); objVal != nil {
				err = objVal.Unmarshal(&obj)
			}
		} else {
			err = objNode.Unmarshal(&obj)
		}
		if err == nil && run != nil {
			c.configs[filename].fn = append(c.configs[filename].fn, run)
		}
	}
	if err != nil {
		err = fmt.Errorf("[%s] %s", filename, err.Error())
	}
	return
}
