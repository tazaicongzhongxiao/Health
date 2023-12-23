package dtime

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type GoTime struct {
	offset   time.Duration
	location *time.Location
	dTime    *time.Time
}

// New 实例 时区 -12 ~ 12 UTC时间zone传13
func New(zone ...int64) *GoTime {
	var offset int64 = 0
	if len(zone) == 1 {
		offset = zone[0]
		if offset == 13 {
			offset = 0
		} else if offset == 0 {
			offset = 8
		}
	}
	location, _ := time.LoadLocation("UTC")
	return &GoTime{
		location: location,
		offset:   time.Duration(offset),
		dTime:    nil,
	}
}

// GetTime 获取时间对象
func (gt *GoTime) GetTime() *GoTime {
	return gt
}

// TimeOffset
// @Description: 设置时区偏差
func (gt *GoTime) TimeOffset(zone int64) {
	gt.offset = time.Duration(zone)
}

// GetLocation 获取时区
func (gt *GoTime) GetLocation() *time.Location {
	return gt.location
}

// GetGoTime 获取时间
func (gt *GoTime) GetGoTime() *time.Time {
	return gt.dTime
}

// NowUnix 获取当前时间戳
func (gt *GoTime) NowUnix() int64 {
	return time.Now().In(gt.location).UTC().Unix()
}

// NowUnixNano 获取当前纳秒级时间戳
func (gt *GoTime) NowUnixNano() int64 {
	return time.Now().In(gt.location).UTC().UnixNano()
}

// SetUnix 时间戳转时间 将在时间戳基础上处理时区
func (gt *GoTime) SetUnix(int64 int64) (t time.Time) {
	return time.Unix(int64, 0).In(gt.location).Add(gt.offset * time.Hour)
}

// LocaleTimeToUTC 时间戳本地时间转UTC或UTC转本地时间 true 本地时间转UTC false UTC转本地时间
func (gt *GoTime) LocaleTimeToUTC(dTime time.Time, utc bool) time.Time {
	if dTime.IsZero() && gt.dTime != nil {
		dTime = *gt.dTime
	}
	if !dTime.IsZero() {
		if utc {
			return dTime.Add(gt.offset * time.Hour)
		} else {
			return dTime.Add(-gt.offset * time.Hour)
		}
	} else {
		return time.Time{}
	}
}

// NowTime 获取当前时间Time(处理时区)
func (gt *GoTime) NowTime() time.Time {
	return time.Now().In(gt.location)
}

// Now 获取当前时间 年-月-日 时:分:秒(处理时区)
func (gt *GoTime) Now() string {
	return gt.LocaleTimeToUTC(gt.NowTime(), true).Format(TT)
}

// GetYmd 获取年月日(处理时区)
func (gt *GoTime) GetYmd() string {
	return gt.LocaleTimeToUTC(gt.NowTime(), true).Format(YMD)
}

// GetHms 获取时分秒(处理时区)
func (gt *GoTime) GetHms() string {
	return gt.LocaleTimeToUTC(gt.NowTime(), true).Format(HMS)
}

// NowStart 获取当天的开始时间, eg: 2018-01-01 00:00:00(处理时区)
func (gt *GoTime) NowStart(int64 ...int64) time.Time {
	now := time.Time{}
	if gt.dTime != nil {
		now = *gt.dTime
	} else if len(int64) == 1 {
		now = time.Unix(int64[0], 0).In(gt.location).Add(gt.offset * time.Hour)
	} else {
		now = gt.NowTime().Add(gt.offset * time.Hour)
	}
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, gt.location)
}

// NowEnd 获取当天的结束时间, eg: 2018-01-01 23:59:59(处理时区)
func (gt *GoTime) NowEnd(int64 ...int64) time.Time {
	now := time.Time{}
	if gt.dTime != nil {
		now = *gt.dTime
	} else if len(int64) == 1 {
		now = time.Unix(int64[0], 0).In(gt.location).Add(gt.offset * time.Hour)
	} else {
		now = gt.NowTime().Add(gt.offset * time.Hour)
	}
	return time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 1e9-1, gt.location)
}

// NowDayStamp 获取相对与当日的时间 00:00:00 开始时间相隔秒数(处理时区)
func (gt *GoTime) NowDayStamp(layout string) int64 {
	now := gt.NowTime()
	theTime, _ := time.ParseInLocation(TT, now.Format(YMD)+" "+layout, gt.location)
	return theTime.Unix() - time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, gt.location).Unix()
}

// ToUnix 2006-01-02 15:04:05 转 时间戳(处理时区) 得到本地时间 utc空UTC时间
func (gt *GoTime) ToUnix(value string, layout string) int64 {
	theTime, _ := time.ParseInLocation(layout, value, gt.location)
	if layout == YMD {
		now := theTime.Add(gt.offset * time.Hour)
		return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, gt.location).Unix()
	}
	return theTime.Add(-gt.offset * time.Hour).Unix()
}

// ToUnixTime 2006-01-02 15:04:05 转 时间(处理时区) 得到本地时间 utc空UTC时间
func (gt *GoTime) ToUnixTime(value string, layout string) *GoTime {
	theTime, _ := time.ParseInLocation(layout, value, gt.location)
	if layout == YMD {
		now := theTime.Add(gt.offset * time.Hour)
		theTime = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, gt.location)
	} else {
		theTime = theTime.Add(-gt.offset * time.Hour)
	}
	gt.dTime = &theTime
	return gt
}

// Timestamp 2006-01-02 15:04:05 转 时间戳(处理时区)
func (gt *GoTime) Timestamp(value string, layout ...string) int64 {
	var _layout string
	if len(layout) > 0 && len(layout[0]) > 0 {
		_layout = layout[0]
	} else {
		_layout = GetDateLayout(value)
	}
	return gt.ToUnix(value, _layout)
}

// Before 当前时间 减去 多少秒
func (gt *GoTime) Before(beforeSecond int64) int64 {
	return gt.NowUnix() - beforeSecond
}

// Next 当前时间 加上 多少秒
func (gt *GoTime) Next(beforeSecond int64) int64 {
	return gt.NowUnix() + beforeSecond
}

// Format time.Time struct to string
// MM - month - 01
// M - month - 1, single bit
// DD - day - 02
// D - day 2
// YYYY - year - 2006
// YY - year - 06
// HH - 24 hours - 03
// H - 24 hours - 3
// hh - 12 hours - 03
// h - 12 hours - 3
// mm - minute - 04
// m - minute - 4
// ss - second - 05
// s - second = 5
func (gt *GoTime) Format(t time.Time, format string) string {
	res := strings.Replace(format, "MM", t.Format("01"), -1)
	res = strings.Replace(res, "M", t.Format("1"), -1)
	res = strings.Replace(res, "DD", t.Format("02"), -1)
	res = strings.Replace(res, "D", t.Format("2"), -1)
	res = strings.Replace(res, "YYYY", t.Format("2006"), -1)
	res = strings.Replace(res, "YY", t.Format("06"), -1)
	res = strings.Replace(res, "HH", fmt.Sprintf("%02d", t.Hour()), -1)
	res = strings.Replace(res, "H", fmt.Sprintf("%d", t.Hour()), -1)
	res = strings.Replace(res, "hh", t.Format("03"), -1)
	res = strings.Replace(res, "h", t.Format("3"), -1)
	res = strings.Replace(res, "mm", t.Format("04"), -1)
	res = strings.Replace(res, "m", t.Format("4"), -1)
	res = strings.Replace(res, "ss", t.Format("05"), -1)
	res = strings.Replace(res, "s", t.Format("5"), -1)
	return res
}

// TimeToHuman 根据时间戳获得人类可读时间
func (gt *GoTime) TimeToHuman(ts int) string {
	var res = ""
	if ts == 0 {
		return res
	}
	t := int(gt.NowUnix()) - ts
	data := [7]map[string]interface{}{
		{"key": 31536000, "value": "年"},
		{"key": 2592000, "value": "个月"},
		{"key": 604800, "value": "星期"},
		{"key": 86400, "value": "天"},
		{"key": 3600, "value": "小时"},
		{"key": 60, "value": "分钟"},
		{"key": 1, "value": "秒"},
	}
	for _, v := range data {
		var c = t / v["key"].(int)
		if 0 != c {
			suffix := "前"
			if c < 0 {
				suffix = "后"
				c = -c
			}
			res = strconv.Itoa(c) + v["value"].(string) + suffix
			break
		}
	}
	return res
}

// Countdown 时间倒计时
func (gt *GoTime) Countdown(end int64) (day, hour, minute, second int64) {
	diff := end - gt.NowUnix()
	if diff <= 0 {
		return
	}
	day = diff / 86400
	diff = diff - day*86400
	hour = diff / 3600
	diff = diff - hour*3600
	minute = diff / 60
	second = diff - minute*60
	return
}

func GetDateLayout(value string) (layout string) {
	layout = "2006-01-02 15:04:05"
	count := len(value)
	switch {
	case count == 10 && strings.Count(value, "-") == 2:
		layout = "2006-01-02"
	case count == 13 && strings.Count(value, "-") == 2:
		layout = "2006-01-02 15"
	case count == 16 && strings.Count(value, "-") == 2:
		layout = "2006-01-02 15:04"
	}
	return
}

// GetDateStart
// @Description: 获取多少天、月、年 之间的时间起始日期时间戳
// @param name day 日 month 月 year 年
// @param d
// @return startDate 开始日期时间戳
// @return endDate 结束日期时间戳
func GetDateStart(name string, d int) (startDate, endDate int64) {
	now := time.Now()
	if name == "day" {
		startDate = time.Date(now.Year(), now.Month(), now.Day()-d, 0, 0, 0, 0, time.UTC).Unix()
		endDate = time.Date(now.Year(), now.Month(), now.Day()-d, 23, 59, 59, 1e9-1, time.UTC).Unix()
	} else if name == "month" {
		startDate = time.Date(now.Year(), now.Month()-time.Month(d), 1, 0, 0, 0, 0, time.UTC).Unix()
		endDate = time.Date(now.Year(), now.Month()-time.Month(d-1), -1, 23, 59, 59, 1e9-1, time.UTC).Unix()
	} else if name == "year" {
		startDate = time.Date(now.Year()-d, 1, 1, 0, 0, 0, 0, time.UTC).Unix()
		endDate = time.Date(now.Year()-d, 12, 31, 23, 59, 59, 1e9-1, time.UTC).Unix()
	}
	return
}
