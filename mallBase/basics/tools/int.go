package tools

import (
	"math"
	"strconv"
)

func IntToString(e int) string {
	return strconv.Itoa(e)
}

func TrunFloat(f float64, prec int) float64 {
	x := math.Pow10(prec)
	return math.Trunc(f*x) / x
}

// ArrayToString 将数字数组转换为字符串
func ArrayToString[T int | int8 | int16 | int32 | int64](str []T) (res []string) {
	res = make([]string, len(str))
	for i := 0; i < len(str); i++ {
		res[i] = strconv.Itoa(int(str[i]))
	}
	return res
}
