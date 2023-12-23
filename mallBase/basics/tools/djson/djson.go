package djson

import (
	json "github.com/json-iterator/go"
)

//
// JsonEncode
// @Description: json编码
//
func JsonEncode(e interface{}) (string, error) {
	if b, err := json.Marshal(e); err == nil {
		return string(b), err
	} else {
		return "", err
	}
}

//
// JsonDecodeMap
// @Description: json解码
//
func JsonDecodeMap(encodeString string) (data map[string]interface{}) {
	_ = json.Unmarshal([]byte(encodeString), &data)
	return data
}

//
// JsonUnmarshal
// @Description: json解码
//
func JsonUnmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, &v)
}
