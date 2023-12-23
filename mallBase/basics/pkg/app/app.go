package app

import (
	"errors"
	"fmt"
	"github.com/TarsCloud/TarsGo/tars"
	"github.com/TarsCloud/TarsGo/tars/util/rogger"
	"github.com/fatih/color"
	"strconv"
	"sync"
	"time"
)

var (
	start    int64
	Cfg      = tars.GetServerConfig()
	CfgName  = "settings.toml"
	services = &Services{services: make(map[string]interface{})}
)

// Services 服务汇总
type Services struct {
	lock     sync.Mutex
	services map[string]interface{}
}

// Logger 获取日志对象
func Logger() *rogger.Logger {
	if Get("logger") != nil {
		return Get("logger").(*rogger.Logger)
	}
	return nil
}

func (service *Services) register(name string, se interface{}) {
	service.lock.Lock()
	defer service.lock.Unlock()
	service.services[name] = se
}

func (service *Services) get(name string) interface{} {
	if val, ok := service.services[name]; ok {
		return val
	}
	return nil
}

// Register 注册其他包的服务
func Register(name string, service interface{}) interface{} {
	services.register(name, service)
	return service
}

// Get 获取其他包的服务
func Get(name string) interface{} {
	return services.get(name)
}

func Err(code int32, format string, args ...interface{}) error {
	if code == 0 {
		code = Fail
	}
	return tars.Errorf(code, fmt.Sprintf(format, args...))
}

// ErrCode 获取 tars错误code 返回1为系统err错误
func ErrCode(err error) int32 {
	return tars.GetErrorCode(err)
}

func IdInt(Id string, empty bool) (ID uint64, err error) {
	if Id != "" {
		ID, errs := strconv.ParseUint(Id, 10, 64)
		if errs == nil {
			return ID, errs
		} else {
			return 0, errors.New("IdIsInt")
		}
	} else if empty == true {
		return 0, errors.New("IdIsInt")
	}
	return 0, nil
}

// Println
// @Description: start 调试一段程序消耗的 时间
// @param c 1Red 2Green 3Yellow 4Blue 5Magenta 6Cyan 7White
// @param req
func Println(c int, req ...interface{}) {
	if start > 0 {
		tmp := time.Now().UnixNano() / 1e6
		color.New(color.Attribute(34)).Println(fmt.Sprintf("结束：%d 耗时：%d 毫秒", tmp, tmp-start))
		start = 0
	}
	value := []string{time.Now().Format("2006-01-02 15:04:05.000")}
	for _, info := range req {
		v, ok := info.(string)
		if ok {
			if start == 0 && v == "start" {
				start = time.Now().UnixNano() / 1e6
				color.New(color.Attribute(34)).Println(fmt.Sprintf("开始：%d", start))
			}
			value = append(value, v)
		} else {
			value = append(value, string(Struct2Json(info)))
		}
	}
	color.New(color.Attribute(c + 30)).Println(value)
	return
}
