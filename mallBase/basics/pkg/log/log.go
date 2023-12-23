package log

import (
	"MyTestMall/mallBase/basics/pkg/app"
	"MyTestMall/mallBase/basics/pkg/config"
	"bytes"
	"fmt"
	"github.com/TarsCloud/TarsGo/tars"
	"github.com/TarsCloud/TarsGo/tars/util/rogger"
	json "github.com/json-iterator/go"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

var (
	isInit    bool
	SLOG      *rogger.Logger
	conf      LogConfig
	logFormat rogger.LogFormat
	logLevel  rogger.LogLevel
)

type LogConfig struct {
	Remote bool `mapstructure:"remote"`
	Hour   int  `mapstructure:"hour"`
}

type JsonLog struct {
	Pre   string      `json:"pre,omitempty"`
	Time  string      `json:"time"`
	Func  string      `json:"func"`
	File  string      `json:"file"`
	Level string      `json:"level"`
	Msg   string      `json:"msg"`
	Data  interface{} `json:"data"`
}

// Start 启动日志服务
func init() {
	if err := config.Config().Bind(app.CfgName, "log", &conf, func() {
		if config.Config().Bind(app.CfgName, "log", &conf, nil) == nil {
			initialize(true)
		}
	}); err != nil {
		fmt.Printf("日志配置读取错误: %s", err.Error())
		panic(err)
		return
	}
	SetLogConfig(conf)
}

func SetLogConfig(m LogConfig) {
	conf = m
	initialize(false)
	//设置日志 等级
	rogger.SetLevel(rogger.StringToLevel(app.Cfg.LogLevel))
	rogger.Colored()
}

func initialize(reentry bool) {
	if conf.Remote && reentry == false {
		SLOG = tars.GetRemoteLogger("TLOG")
		SLOG.Info("启动远程日志...")
	} else if conf.Remote == false {
		SLOG = tars.GetLogger("TLOG")
		SLOG.Info("启动本地日志...")
	}
	if conf.Hour == 0 {
		_ = SLOG.SetFileRoller(app.Cfg.LogPath, 1, int(app.Cfg.LogSize))
	} else if conf.Hour%24 == 0 {
		_ = SLOG.SetDayRoller(app.Cfg.LogPath, conf.Hour/24)
	} else if conf.Hour > 0 {
		_ = SLOG.SetHourRoller(app.Cfg.LogPath, conf.Hour)
	}
	if app.Cfg.SampleType == "json" {
		rogger.SetFormat(rogger.Json)
	} else {
		rogger.SetFormat(rogger.Text)
	}
	logFormat = rogger.GetLogFormat()
	logLevel = rogger.GetLogLevel()
	SLOG.SetConsole()
	app.Register("logger", SLOG)
	isInit = true
}

func Debug(msg string, data interface{}, depth ...int) {
	if rogger.DEBUG < logLevel {
		return
	}
	d := WriteF(rogger.DEBUG, msg, data, depth...)
	if isInit {
		SLOG.WriteLog(d)
	} else {
		fmt.Println("DEBUG | ", string(d))
	}
	return
}

func Info(msg string, data interface{}, depth ...int) {
	if rogger.INFO < logLevel {
		return
	}
	d := WriteF(rogger.INFO, msg, data, depth...)
	if isInit {
		SLOG.WriteLog(d)
	} else {
		fmt.Println("INFO | ", string(d))
	}
	return
}

func Warn(msg string, data interface{}, depth ...int) {
	if rogger.WARN < logLevel {
		return
	}
	d := WriteF(rogger.WARN, msg, data, depth...)
	if isInit {
		SLOG.WriteLog(d)
	} else {
		fmt.Println("WARN | ", string(d))
	}
	return
}

func Error(msg string, data interface{}, depth ...int) {
	if rogger.INFO < logLevel {
		return
	}
	d := WriteF(rogger.ERROR, msg, data, depth...)
	if isInit {
		SLOG.WriteLog(d)
	} else {
		fmt.Println("ERROR | ", string(d))
	}
	return
}

func FuncList() {
	var stacktrace string
	for i := 1; ; i++ {
		pc, file, line, got := runtime.Caller(i)
		if !got {
			break
		}
		file = filepath.Base(file)
		Func := getFuncName(runtime.FuncForPC(pc).Name())
		stacktrace += fmt.Sprintf("%s %s:%d %d\n", Func, file, line, i)
	}
	fmt.Println(fmt.Sprintf("\n%s", stacktrace))
	return
}

func WriteF(level rogger.LogLevel, msg string, data interface{}, depth ...int) []byte {
	skip := 2
	if len(depth) == 1 {
		skip = depth[0]
	}
	if logFormat == rogger.Json {
		log := JsonLog{}
		log.Time = time.Now().Format("2006-01-02 15:04:05.000")
		pc, file, line, ok := runtime.Caller(skip)
		if !ok {
			file = "???"
			line = 0
		} else {
			file = filepath.Base(file)
		}
		log.Func = getFuncName(runtime.FuncForPC(pc).Name())
		log.File = fmt.Sprintf("%s:%d", file, line)
		log.Level = level.String()
		log.Msg = msg
		log.Data = data
		buf := &bytes.Buffer{}
		encoder := json.NewEncoder(buf)
		encoder.SetEscapeHTML(false)
		_ = encoder.Encode(log)
		return buf.Bytes()
	} else {
		buf := bytes.NewBuffer(nil)
		fmt.Fprintf(buf, "%s|", time.Now().Format("2006-01-02 15:04:05.000"))
		pc, file, line, ok := runtime.Caller(skip)
		if !ok {
			file = "???"
			line = 0
		} else {
			file = filepath.Base(file)
		}
		fmt.Fprintf(buf, "%s:%s:%d|", file, getFuncName(runtime.FuncForPC(pc).Name()), line)
		buf.WriteString(coloredString(level))
		buf.WriteByte('|')
		fmt.Fprint(buf, msg)
		buf.WriteByte('|')
		v, ok := data.(string)
		if ok {
			fmt.Fprint(buf, v)
		} else {
			jsons, _ := json.Marshal(data)
			fmt.Fprint(buf, string(jsons))
		}
		buf.WriteByte('\n')
		return buf.Bytes()
	}
}

// Colored enable colored level string when use console writer
func coloredString(level rogger.LogLevel) string {
	switch level {
	case rogger.DEBUG:
		return "\x1b[34mDEBUG\x1b[0m" // blue
	case rogger.INFO:
		return "\x1b[32mINFO\x1b[0m" //green
	case rogger.WARN:
		return "\x1b[33mWARN\x1b[0m" // yellow
	case rogger.ERROR:
		return "\x1b[31mERROR\x1b[0m" //cred
	default:
		return "\x1b[37mUNKNOWN\x1b[0m" // white
	}
}

func getFuncName(name string) string {
	idx := strings.LastIndexByte(name, '/')
	if idx != -1 {
		name = name[idx:]
		idx = strings.IndexByte(name, '.')
		if idx != -1 {
			name = strings.TrimPrefix(name[idx:], ".")
		}
	}
	return name
}
