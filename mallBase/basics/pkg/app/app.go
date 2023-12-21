package app

import (
	"fmt"
	"github.com/TarsCloud/TarsGo/tars"
	"github.com/fatih/color"
	"time"
)

var (
	start   int64
	Cfg     = tars.GetServerConfig()
	CfgName = "settings.toml"
)

// Println
// @Description: start 调试一段程序消耗的 时间
// @param c 1Red 2Green 3Yellow 4Blue 5Magenta 6Cyan 7White
// @param req
func Println(c int, req ...interface{}) {
	if start > 0 {
		tmp := time.Now().UnixNano() / 1e6
		color.New(color.Attribute(34)).Println(fmt.Sprintf("结束：%d 耗时 %d 毫秒", tmp, tmp-start))
		start = 0
	}
	value := []string{time.Now().Format("2006-01-02 15:04:05.000")}
	for _, info := range req {
		if v, ok := info.(string); ok {
			if start == 0 && v == "start" {
				start = time.Now().UnixNano() / 1e6
				color.New(color.Attribute(34)).Println(fmt.Sprintf("开始：%d", start))
			}
			value = append(value, v)
		} else {
			value = append(value, string(Struct2Json(v)))
		}

	}
	color.New(color.Attribute(c + 30)).Println(value)
	return
}
