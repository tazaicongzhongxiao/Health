package log

import (
	"strings"
)

type NsqLoggerI struct {
}

func (l *NsqLoggerI) Output(calldepth int, s string) error {
	level := strings.Split(s, " ")[0]
	switch level {
	case "INF":
		Info(s, nil)
	case "WRN":
		Warn(s, nil)
	case "ERR":
		Error(s, nil)
	default:
		Debug(s, nil)
	}
	return nil
}

func NsqLogger() *NsqLoggerI {
	return &NsqLoggerI{}
}
