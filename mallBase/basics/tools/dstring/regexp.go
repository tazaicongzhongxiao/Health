package dstring

import (
	"fmt"
	"regexp"
)

func MatchString(pattern string, str string) bool {
	ok, err := regexp.MatchString(pattern, str)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return ok
}

// IsNumber 是数字
// @param length 长度验证： 1个值时为指定长度；2个值时分别为 min 和 max
func IsNumber(data string, length ...uint) bool {
	if len(length) == 1 && length[0] > 0 {
		return MatchString(fmt.Sprintf(`^[0-9]{%d}$`, length[0]), data)
	} else if len(length) == 2 && length[0] <= length[1] && length[0] > 0 {
		return MatchString(fmt.Sprintf(`^[0-9]{%d,%d}$`, length[0], length[1]), data)
	} else {
		return MatchString(`^[0-9]+$`, data)
	}
}

// HasNumber 包含数字
func HasNumber(data string) bool {
	return MatchString(`[0-9]`, data)
}

// IsDecimal 是实数
// @param scale 小数位长度验证： 1个值时为指定长度；2个值时分别为 min 和 max
func IsDecimal(data string, scale ...uint) bool {
	if len(scale) == 1 {
		if scale[0] == 0 {
			return MatchString(`^-?(([1-9]\d*)|0)$`, data)
		} else {
			return MatchString(fmt.Sprintf(`^-?(([1-9]\d*)|0)[.]\d{%d}$`, scale[0]), data)
		}
	} else if len(scale) == 2 && scale[0] <= scale[1] {
		if scale[0] == 0 && scale[1] == 0 {
			return MatchString(`^-?(([1-9]\d*)|0)$`, data)
		} else if scale[0] == 0 && scale[1] > 0 {
			return MatchString(`^-?(([1-9]\d*)|0)$`, data) || MatchString(fmt.Sprintf(`^-?(([1-9]\d*)|0)([.]\d{1,%d})?$`, scale[1]), data)
		} else {
			return MatchString(fmt.Sprintf(`^-?(([1-9]\d*)|0)[.]\d{%d,%d}$`, scale[0], scale[1]), data)
		}
	} else {
		return MatchString(`^-?(([1-9]\d*)|0)([.]\d+)?$`, data)
	}
}

// IsUDecimal 是非负实数
// @param scale 小数位长度验证： 1个值时为指定长度；2个值时分别为 min 和 max
func IsUDecimal(data string, scale ...uint) bool {
	if len(scale) == 1 {
		if scale[0] == 0 {
			return MatchString(`^(([1-9]\d*)|0)$`, data)
		} else {
			return MatchString(fmt.Sprintf(`^(([1-9]\d*)|0)[.]\d{%d}$`, scale[0]), data)
		}
	} else if len(scale) >= 2 && scale[0] <= scale[1] {
		if scale[0] == 0 && scale[1] == 0 {
			return MatchString(`^(([1-9]\d*)|0)$`, data)
		} else if scale[0] == 0 && scale[1] > 0 {
			return MatchString(`^(([1-9]\d*)|0)$`, data) || MatchString(fmt.Sprintf(`^(([1-9]\d*)|0)([.]\d{1,%d})?$`, scale[1]), data)
		} else {
			return MatchString(fmt.Sprintf(`^(([1-9]\d*)|0)[.]\d{%d,%d}$`, scale[0], scale[1]), data)
		}
	} else {
		return MatchString(`^(([1-9]\d*)|0)([.]\d+)?$`, data)
	}
}

// IsInt 是整数
func IsInt(data string) bool {
	return MatchString(`^((-?([1-9]\d*))|0)$`, data)
}

// IsUInt 是非负整数
func IsUInt(data string) bool {
	return MatchString(`^(([1-9]\d*)|0)$`, data)
}

// IsAlpha 是全英文字母
// @param length 长度验证： 1个值时为指定长度；2个值时分别为 min 和 max
func IsAlpha(data string, length ...uint) bool {
	if len(length) == 1 && length[0] > 0 {
		return MatchString(fmt.Sprintf(`^[A-Za-z]{%d}$`, length[0]), data)
	} else if len(length) == 2 && length[0] <= length[1] && length[0] > 0 {
		return MatchString(fmt.Sprintf(`^[A-Za-z]{%d,%d}$`, length[0], length[1]), data)
	} else {
		return MatchString("^[A-Za-z]+$", data)
	}
}

// HasAlpha 包含英文字母
func HasAlpha(data string) bool {
	return MatchString("[A-Za-z]", data)
}

// IsUpper 是大写字母
func IsUpper(char string) bool {
	return MatchString("^[A-Z]$", char)
}

// HasUpper 包含大写字母
func HasUpper(str string) bool {
	return MatchString("[A-Z]", str)
}

// IsLower 是小写字母
func IsLower(char string) bool {
	return MatchString("^[a-z]$", char)
}

// HasLower 包含小写字母
func HasLower(str string) bool {
	return MatchString("[a-z]", str)
}

// HasChar 包含标点字符
func HasChar(data string) bool {
	return MatchString(`[\.~!@#$%^&*()\-=_+:;,?]`, data)
}

// IsChinese 是全中文汉字
func IsChinese(data string) bool {
	return MatchString("^\\p{Han}+$", data)
	//return MatchString("^[\u4e00-\u9fa5]+$", data)
}

// HasChinese 包含中文汉字
func HasChinese(data string) bool {
	return MatchString("\\p{Han}", data)
	//return MatchString("[\u4e00-\u9fa5]", data)
}

// IsNumAlpha 是数字字母
// @param length 长度验证： 1个值时为指定长度；2个值时分别为 min 和 max
func IsNumAlpha(data string, length ...uint) bool {
	if len(length) == 1 && length[0] > 0 {
		return MatchString(fmt.Sprintf(`^[0-9A-Za-z]{%d}$`, length[0]), data)
	} else if len(length) == 2 && length[0] <= length[1] && length[0] > 0 {
		return MatchString(fmt.Sprintf(`^[0-9A-Za-z]{%d,%d}$`, length[0], length[1]), data)
	} else {
		return MatchString("^[0-9A-Za-z]+$", data)
	}
}

// HasNumAlpha 包含数字字母
func HasNumAlpha(data string) bool {
	return HasNumber(data) && HasAlpha(data)
}

// IsNum_Alpha 是数字字母下划线
// @param length 长度验证： 1个值时为指定长度；2个值时分别为 min 和 max
func IsNum_Alpha(data string, length ...uint) bool {
	if len(length) == 1 && length[0] > 0 {
		return MatchString(fmt.Sprintf(`^[0-9A-Za-z_]{%d}$`, length[0]), data)
	} else if len(length) == 2 && length[0] <= length[1] && length[0] > 0 {
		return MatchString(fmt.Sprintf(`^[0-9A-Za-z_]{%d,%d}$`, length[0], length[1]), data)
	} else {
		return MatchString("^[0-9A-Za-z_]+$", data)
	}
}

// HasNum_Alpha 包含数字字母下划线
func HasNum_Alpha(data string) bool {
	return HasNumber(data) && HasAlpha(data) && MatchString("[_]", data)
}
