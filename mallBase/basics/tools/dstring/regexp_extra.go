package dstring

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

/*
公民身份证号
xxxxxx yyyy MM dd 375 0     十八位
xxxxxx   yy MM dd  75 0     十五位
地区：[1-9]\d{5}
年的前两位：(18|19|([23]\d))      1800-2399
年的后两位：\d{2}
月份：((0[1-9])|(10|11|12))
天数：(([0-2][1-9])|10|20|30|31) 闰年不能禁止29+
三位顺序码：\d{3}
两位顺序码：\d{2}
校验码：   [0-9Xx]
十八位：^[1-9]\d{5}(18|19|([23]\d))\d{2}((0[1-9])|(10|11|12))(([0-2][1-9])|10|20|30|31)\d{3}[0-9Xx]$
十五位：^[1-9]\d{5}\d{2}((0[1-9])|(10|11|12))(([0-2][1-9])|10|20|30|31)\d{3}$
总：
(^[1-9]\d{5}(18|19|([23]\d))\d{2}((0[1-9])|(10|11|12))(([0-2][1-9])|10|20|30|31)\d{3}[0-9Xx]$)|(^[1-9]\d{5}\d{2}((0[1-9])|(10|11|12))(([0-2][1-9])|10|20|30|31)\d{3}$)
*/
func IsIDCard(data string) bool {
	if len(data) == 18 {
		return checkIDCardLast(data) && isIDCard(data)
	} else {
		return isIDCard(data)
	}
}

func isIDCard(data string) bool {
	return MatchString(`(^[1-9]\d{5}(18|19|([23]\d))\d{2}((0[1-9])|(10|11|12))(([0-2][1-9])|10|20|30|31)\d{3}[0-9Xx]$)|(^[1-9]\d{5}\d{2}((0[1-9])|(10|11|12))(([0-2][1-9])|10|20|30|31)\d{2}$)`, data)
}

// 第二代身份证校验码（GB 11643-1999）
func checkIDCardLast(data string) bool {
	cardNo := strings.ToUpper(data)
	checks := []string{"1", "0", "X", "9", "8", "7", "6", "5", "4", "3", "2"}
	var sum int
	for i := 17; i > 0; i-- {
		n, _ := strconv.Atoi(cardNo[17-i : 18-i])
		p := int(math.Pow(2, float64(i))) % 11
		sum += n * p
	}
	return cardNo[17:] == checks[sum%11]
}

// checkLuhn 银行卡号校验：Luhn算法（模10算法）
// Check 检查银行卡号是否合法
func IsBankCard(cardNumber string) bool {
	sum := 0
	for index := len(cardNumber) - 1; index >= 0; index-- {
		if !(cardNumber[index] >= '0' && cardNumber[index] <= '9') {
			return false
		}
		t := int(cardNumber[index] - '0')
		if (len(cardNumber)-index)%2 == 0 {
			// 偶数位
			n := t * 2
			for n > 0 {
				sum += n % 10
				n /= 10
			}
		} else {
			// 奇数位
			sum += t
		}
	}
	return sum%10 == 0
}

/*
验证所给手机号码是否符合手机号的格式.
移动: 134、135、136、137、138、139、150、151、152、157、158、159、182、183、184、187、188、178(4G)、147(上网卡)；
联通: 130、131、132、155、156、185、186、176(4G)、145(上网卡)、175；
电信: 133、153、180、181、189 、177(4G)；
卫星通信:  1349
虚拟运营商: 170、173
2018新增: 16x, 19x
*/
func IsMobile(data string) bool {
	return MatchString(`^13[\d]{9}$|^14[5,7]{1}\d{8}$|^15[^4]{1}\d{8}$|^16[\d]{9}$|^17[0,3,5,6,7,8]{1}\d{8}$|^18[\d]{9}$|^19[\d]{9}$`, data)
}

// IsTel 国内座机电话号码："XXXX-XXXXXXX"、"XXXX-XXXXXXXX"、"XXX-XXXXXXX"、"XXX-XXXXXXXX"、"XXXXXXX"、"XXXXXXXX"
func IsTel(data string) bool {
	return MatchString(`^((\d{3,4})|\d{3,4}-)?\d{7,8}$`, data)
}

// IsEmail Email地址
func IsEmail(data string) bool {
	return MatchString(`^[a-zA-Z0-9_\-\.]+@[a-zA-Z0-9_\-]+(\.[a-zA-Z0-9_\-]+)+$`, data)
}

// IsIPAddress 是否IPv4地址
func IsIPAddress(ip string) bool {
	return MatchString("(25[0-5]|2[0-4]\\d|[0-1]\\d{2}|[1-9]?\\d)\\.(25[0-5]|2[0-4]\\d|[0-1]\\d{2}|[1-9]?\\d)\\.(25[0-5]|2[0-4]\\d|[0-1]\\d{2}|[1-9]?\\d)\\.(25[0-5]|2[0-4]\\d|[0-1]\\d{2}|[1-9]?\\d)", ip)
}

// IsIntranetIP 是否内网IP地址
func IsIntranetIP(s string) bool {
	return MatchString(`^((192\.168|172\.([1][6-9]|[2]\d|3[01]))(\.([2][0-4]\d|[2][5][0-5]|[01]?\d?\d)){2}|10(\.([2][0-4]\d|[2][5][0-5]|[01]?\d?\d)){3})$`, s)
}

// IsURL URL地址
func IsURL(data string) bool {
	return MatchString(`(https?|ftp|file)://[-A-Za-z0-9+&@#/%?=~_|!:,.;]+[-A-Za-z0-9+&@#/%=~_|]`, data)
}

// IsGrpc 检测GRPC
func IsGrpc(data string) bool {
	return MatchString(`^[A-Za-z]+\.[A-Za-z]+\.[A-Za-z]+$|^[A-Za-z]+\.[A-Za-z]+\.[A-Za-z]+@[0-9A-Za-z\s-.]+$`, data)
}

// IsNSQ 检测NSQ 地址
func IsNSQ(data string) bool {
	return MatchString(`^[A-Za-z]+\.[A-Za-z]+\.[A-Za-z]+$`, data)
}

// IsMac Mac地址
func IsMac(data string) bool {
	return MatchString(`^([0-9A-Fa-f]{2}[\-:]){5}[0-9A-Fa-f]{2}$`, data)
}

// IsQQ 腾讯QQ号，从10000开始
func IsQQ(data string) bool {
	return MatchString(`^[1-9][0-9]{4,}$`, data)
}

// IsPostCode 邮政编码
func IsPostCode(data string) bool {
	return MatchString(`^\d{6}$`, data)
}

// IsDateFormat sep 1日期 2时间 3日期加时间
func IsDateFormat(data string, sep int) bool {
	var pattern string
	if sep == 1 { // 日期 2014-01-01
		pattern = `^[1-9]\d{3}-(0[1-9]|1[0-2])-(0[1-9]|[1-2][0-9]|3[0-1])$`
	} else if sep == 2 { // 时间 12:00:00
		pattern = `^(20|21|22|23|[0-1]\d):[0-5]\d:[0-5]\d$`
	} else { // 日期加时间 2014-01-01 12:00:00
		pattern = `^[1-9]\d{3}-(0[1-9]|1[0-2])-(0[1-9]|[1-2][0-9]|3[0-1])\s+(20|21|22|23|[0-1]\d):[0-5]\d:[0-5]\d$`
	}
	return MatchString(pattern, data)
}

// IsAccount
// @Description: 检查账号（字母开头，数字字母下划线）
// length 长度验证： 1个值时为指定长度；2个值时分别为 min 和 max
func IsAccount(data string, length ...uint) bool {
	var lengthStr string
	if len(length) == 1 && length[0] > 0 {
		lengthStr = fmt.Sprintf("{%d}", length[0]-1)
	} else if len(length) == 2 && length[0] <= length[1] && length[0] > 0 {
		lengthStr = fmt.Sprintf("{%d,%d}", length[0]-1, length[1]-1)
	} else {
		lengthStr = "{5,19}"
	}
	return MatchString(fmt.Sprintf(`^[A-Za-z]{1}[0-9A-Za-z_]%s$`, lengthStr), data)
}

// IsPwd 检查密码
// @param length 长度验证： 1个值时为指定长度；2个值时分别为 min 和 max
func IsPwd(data string, level uint, length ...uint) bool {
	var lengthStr string
	if len(length) == 1 && length[0] > 0 {
		lengthStr = fmt.Sprintf("{%d}", length[0])
	} else if len(length) == 2 && length[0] <= length[1] && length[0] > 0 {
		lengthStr = fmt.Sprintf("{%d,%d}", length[0], length[1])
	} else {
		lengthStr = "{6,20}"
	}
	switch level {
	case 1: //包含数字、字母
		return MatchString(fmt.Sprintf(`^[\w\S]%s$`, lengthStr), data) && HasNumber(data) && HasAlpha(data)
	case 2: //包含数字、字母、下划线
		return MatchString(fmt.Sprintf(`^[\w\S]%s$`, lengthStr), data) && HasNum_Alpha(data)
	case 3: //包含数字、字母、特殊字符
		return MatchString(fmt.Sprintf(`^[\w\S]%s$`, lengthStr), data) && HasNumber(data) && HasAlpha(data) && HasChar(data)
	case 4: //包含数字、大小写字母
		return MatchString(fmt.Sprintf(`^[\w\S]%s$`, lengthStr), data) && HasNumber(data) && HasUpper(data) && HasLower(data)
	case 5: //包含数字、大小写字母、下划线
		return MatchString(fmt.Sprintf(`^[\w\S]%s$`, lengthStr), data) && HasNumber(data) && HasUpper(data) && HasLower(data) && MatchString("[_]", data)
	case 6: //包含数字、大小写字母、特殊字符
		return MatchString(fmt.Sprintf(`^[\w\S]%s$`, lengthStr), data) && HasNumber(data) && HasUpper(data) && HasLower(data) && HasChar(data)
	default:
		return MatchString(fmt.Sprintf(`^[\w\S]%s$`, lengthStr), data)
	}
}
