package dstring

import "strings"

// SubstrByEnd
// 截取字符串
// end 0：截取全部；负数：从后往前
func SubstrByEnd(s string, start int, end int) string {
	rs := []rune(s)
	if start < 0 {
		start = 0
	}
	if start > len(rs) {
		start = start % len(rs)
	}
	if end >= 0 {
		if end < start || end > len(rs) {
			end = len(rs)
		}
	} else {
		if len(rs)+end < start {
			end = len(rs)
		} else {
			end = len(rs) + end
		}
	}
	return string(rs[start:end])
}

// HideNo
// 隐藏字符串
// start：前端显示长度
// end：后端显示长度
// length：指定显示总长度，若不指定，则按原字符串长度输出
func HideNo(s string, start int, end int, length ...int) string {
	s = strings.TrimSpace(s)
	oldLen := len(s)
	newLen := oldLen
	if len(length) > 0 {
		newLen = length[0]
	}
	minLen := oldLen
	if oldLen >= newLen {
		minLen = newLen
	}
	if minLen <= 1 {
		return strings.Repeat("*", newLen)
	}
	if start >= minLen {
		start = minLen - 1
		end = 0
	} else if end >= minLen {
		start = 0
		end = minLen - 1
	} else if start+end >= minLen {
		start = minLen / 2
		end = minLen/2 - 1
	}
	rs := Substr(s, 0, start) + strings.Repeat("*", newLen-start-end) + Substr(s, 0, -end)
	return rs
}

// HidePhone
// 隐藏 手机号
func HidePhone(s string) string {
	s = strings.TrimSpace(s)
	length := len(s)
	if length == 0 {
		return ""
	}
	if strings.Contains(s, "+") {
		return Substr(s, 0, length-8) + "****" + SubstrByEnd(s, length-4, 0)
	} else {
		if strings.Contains(s, "-") || strings.Contains(s, "_") || strings.Contains(s, " ") {
			return Substr(s, 0, length-6) + "***" + SubstrByEnd(s, length-3, 0)
		} else {
			if length == 11 {
				return Substr(s, 0, 3) + "****" + SubstrByEnd(s, length-4, 0)
			} else {
				return Substr(s, 0, length-6) + "***" + SubstrByEnd(s, length-3, 0)
			}
		}
	}
}

// HideEmail
// 隐藏 邮箱
func HideEmail(s string) string {
	emails := strings.Split(s, "@")
	if len(emails) != 2 {
		return s
	}
	return HideNo(emails[0], 2, 2, 6) + "@" + emails[1]
}

// HidePwd
// 隐藏 密码
func HidePwd(s string, allHide ...bool) string {
	s = strings.TrimSpace(s)
	if len(allHide) > 0 && allHide[0] {
		return "******"
	} else {
		if len(s) > 0 {
			return "******"
		} else {
			return ""
		}
	}
}
