package dstring

import (
	"regexp"
	"sort"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

// TrimRightSpace
// @Description: 去掉换行,空格,回车,制表
// @param s
// @return string
func TrimRightSpace(s string) string {
	return strings.TrimRight(s, "\r\n\t ")
}

// StringToInt64
// @Description: 字符串转int64
// @param e
// @return res
func StringToInt64(e string) (res int64) {
	res, _ = strconv.ParseInt(e, 10, 64)
	return res
}

// StringToUint64
// @Description: 字符串转uint64
// @param e
// @return res
func StringToUint64(e string) (res uint64) {
	res, _ = strconv.ParseUint(e, 10, 64)
	return res
}

func StringFloatToInt(e string, prec int) (int, error) {
	s, _ := strconv.ParseFloat(e, 64)
	return strconv.Atoi(strconv.FormatFloat(s, 'f', prec, 64))
}

// StringToFloat64
// @Description: 字符串转float64
// @param e
// @return float64
func StringToFloat64(e string) float64 {
	s, _ := strconv.ParseFloat(e, 64)
	return s
}

// StringToBool
// @Description: 字符串转bool
// @param e
// @return bool
func StringToBool(e string) bool {
	s, _ := strconv.ParseBool(e)
	return s
}

// StringToInt
// @Description: 字符串转int
// @param e
// @return int
func StringToInt(e string) int {
	s, _ := strconv.Atoi(e)
	return s
}

// Substr
// @Description: 截取字符串
// @param s
// @param start
// @param length 不设置：截取全部；负数：向前截取
// @return string
func Substr(s string, start int, length ...int) string {
	rs := []rune(s)
	l := len(rs)
	if len(length) > 0 {
		l = length[0]
	}
	if l > 0 {
		if start <= 0 {
			start = 0
		} else {
			if start > len(rs) {
				start = start % len(rs)
			}
		}
		end := start + l
		if start+l > len(rs) {
			end = len(rs)
		}
		return string(rs[start:end])
	} else if l < 0 {
		if start <= 0 {
			start = len(rs)
		} else {
			if start > len(rs) {
				start = start % len(rs)
			}
		}
		end := start
		start = end + l
		if end+l < 0 {
			start = 0
		}
		return string(rs[start:end])
	} else {
		return ""
	}
}

// StrPos
// @Description:描述：查找字符串在另一个字符串中，第一次出现的位置，没有找到子串，则返回-1；注意：对大小写敏感；
// @param sub 要查找的字符串
// @param str 搜索的字符串
// @return int
func StrPos(sub, str string) int {
	strList := []rune(str)
	subPos := 0
	isFind := false
	for k := range strList {
		if string(strList[k]) == sub {
			subPos = k
			isFind = true
			break
		}
	}
	if !isFind {
		return -1
	}
	return subPos
}

// StrLastPos
// @Description: 查找字符串在另一个字符串中，最后一次出现的位置, 没有找到子串，则返回-1 注意：对大小写敏感
// @param sub 要查找的字符串
// @param str 搜索的字符串
// @return int
func StrLastPos(sub, str string) int {
	strList := []rune(str)
	subPos := 0
	isFind := false
	for k := range strList {
		if string(strList[k]) == sub {
			subPos = k
			isFind = true
		}
	}
	if !isFind {
		return -1
	}
	return subPos
}

// StrReverse
// @Description: 字符串反转
// @param str 要反转的字符串
// @return string
func StrReverse(str string) string {
	strList := []rune(str)
	strLen := len(strList)
	newStrList := make([]rune, 0, strLen)
	for i := strLen - 1; i >= 0; i-- {
		newStrList = append(newStrList, strList[i])
	}
	return string(newStrList)
}

// ReplaceNoCase
// @Description: 替换字符串（不区分大小写）
// @param s
// @param old
// @param new
// @param n
// @return string
func ReplaceNoCase(s string, old string, new string, n int) string {
	if n == 0 {
		return s
	}
	ls := strings.ToLower(s)
	lold := strings.ToLower(old)
	if m := strings.Count(ls, lold); m == 0 {
		return s
	} else if n < 0 || m < n {
		n = m
	}
	ns := make([]byte, len(s)+n*(len(new)-len(old)))
	w := 0
	start := 0
	for i := 0; i < n; i++ {
		j := start
		if len(old) == 0 {
			if i > 0 {
				_, wid := utf8.DecodeRuneInString(s[start:])
				j += wid
			}
		} else {
			j += strings.Index(ls[start:], lold)
		}
		w += copy(ns[w:], s[start:j])
		w += copy(ns[w:], new)
		start = j + len(old)
	}
	w += copy(ns[w:], s[start:])
	return string(ns[0:w])
}

// CompareStringSlice
// @Description: 对比两个string切片，看内容是否一样
// @param a
// @param b
// @return bool
func CompareStringSlice(a []string, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	sort.Strings(a)
	sort.Strings(b)
	for i, item := range a {
		if item != b[i] {
			return false
		}
	}
	return true
}

// IsAllChinese
// @Description: 是否全是汉字
// @param str
// @return bool
func IsAllChinese(str string) bool {
	if len(str) <= 0 {
		return false
	}
	ret := true
	for _, r := range str {
		if !unicode.Is(unicode.Scripts["Han"], r) {
			ret = false
			break
		}
	}
	return ret
}

// IsNormalStr
// @Description: 是否为汉字、字母、数字
// @param str
// @return bool
func IsNormalStr(str string) bool {
	if len(str) <= 0 {
		return false
	}
	reg := regexp.MustCompile("^[a-zA-Z0-9\u4e00-\u9fa5]+$")
	return reg.MatchString(str)
}

// ToFirstUpper
// @Description: 转换成 首字母大写
// @param s
// @return string
func ToFirstUpper(s string) string {
	s = strings.TrimSpace(s)
	if s != "" {
		s = strings.ToUpper(s[:1]) + s[1:]
	}
	return s
}

// ToFirstLower
// @Description: 转换成 首字母小写
// @param s
// @return string
func ToFirstLower(s string) string {
	s = strings.TrimSpace(s)
	if s != "" {
		s = strings.ToLower(s[:1]) + s[1:]
	}
	return s
}

// ToCamelCase
// @Description: 转换成 大驼峰命名（UserId）
// @param s
// @return string
func ToCamelCase(s string) string {
	if IsNum_Alpha(s) {
		var rs string
		s = strings.TrimSpace(s)
		es := strings.Split(s, "_")
		for _, e := range es {
			rs += ToFirstUpper(e)
		}
		return rs
	} else {
		return s
	}
}

// TocamelCase
// @Description: 转换成 小驼峰命名（userId）
// @param s
// @return string
func TocamelCase(s string) string {
	return ToFirstLower(ToCamelCase(s))
}

// ToUnderscoreCase
// @Description: 转换成 大下划线命名（USER_ID）
// @param s
// @return string
func ToUnderscoreCase(s string) string {
	return strings.ToUpper(TounderscoreCase(s))
}

// TounderscoreCase
// @Description: 转换成 小下划线命名（user_id）
// @param s
// @return string
func TounderscoreCase(s string) string {
	if IsNum_Alpha(s) {
		var rs string
		l := len(s)
		for i := 0; i < l; i++ {
			e := s[i : i+1]
			if IsUpper(e) {
				e = "_" + strings.ToLower(e)
			}
			rs += e
		}
		rs = strings.TrimPrefix(rs, "_")
		rs = strings.Replace(rs, "__", "_", -1)
		return rs
	} else {
		return s
	}
}

// StrPad
// input string 原字符串
// padLength int 规定补齐后的字符串位数
// padString string 自定义填充字符串
// padType string 填充类型:LEFT(向左填充,自动补齐位数), 默认右侧
func StrPad(input string, padLength int, padString string, padType string) string {
	output := ""
	inputLen := len(input)
	if inputLen >= padLength {
		return input
	}
	padStringLen := len(padString)
	needFillLen := padLength - inputLen
	if diffLen := padStringLen - needFillLen; diffLen > 0 {
		padString = padString[diffLen:]
	}
	for i := 1; i <= needFillLen; i += padStringLen {
		output += padString
	}
	switch padType {
	case "LEFT":
		return output + input
	default:
		return input + output
	}
}
