package hbase64

import "encoding/base64"

/**
 * @brief base64编码
 * @param string name 待编码的字符
 * @return mixed
 */
func Base64Encode(data string) string {
	encodeString := base64.StdEncoding.EncodeToString([]byte(data))
	return encodeString
}

/**
 * @brief base64解码
 * @param string encodeString base64编码的字符
 * @return mixed
 */
func Base64Decode(encodeString string) string {
	decodeBytes, err := base64.StdEncoding.DecodeString(encodeString)
	if err == nil {
		return string(decodeBytes)
	} else {
		return ""
	}
}