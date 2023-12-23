package encrypt

// MergeNum 合并2个数字为一个新长数字
func MergeNum(left, right int64) int64 {
	if left < 0 || right < 0 {
		return 0
	}
	var result = left
	result <<= 32
	result |= right
	return result
}

// SplitNum 还原合并的数字
func SplitNum(val int64) (int64, int64) {
	var left = val >> 32
	var right = (val << 32) >> 32
	return left, right
}
