package dinterface

import (
	json "github.com/json-iterator/go"
	"reflect"
	"strconv"
)

// IsNil 判断 interface 值是否为空
func IsNil(value interface{}) bool {
	switch value.(type) {
	case int32:
		return value.(int32) == 0
	case int64:
		return value.(int64) == 0
	case float32:
		return value.(float32) == 0
	case float64:
		return value.(float64) == 0
	case string:
		return value.(string) == ""
	case bool:
		return value.(bool)
	}
	return false
}

// IsVal 验证接口格式  kind 1 数字 2 字符串 3 布尔值 4 interface
func IsVal(value interface{}, kind int8) bool {
	switch reflect.TypeOf(value).Kind() {
	case reflect.Float32, reflect.Float64, reflect.Int8, reflect.Int32, reflect.Int64:
		return kind == 1
	case reflect.String:
		return kind == 2
	case reflect.Bool:
		return kind == 3
	case reflect.Interface, reflect.Ptr:
		return kind == 4
	}
	return false
}

func ConvertAnyToStr(v interface{}) string {
	if v == nil {
		return ""
	}
	switch d := v.(type) {
	case string:
		return d
	case int, int8, int16, int32, int64:
		return strconv.FormatInt(reflect.ValueOf(v).Int(), 10)
	case uint, uint8, uint16, uint32, uint64:
		return strconv.FormatUint(reflect.ValueOf(v).Uint(), 10)
	case []byte:
		return string(d)
	case float32, float64:
		return strconv.FormatFloat(reflect.ValueOf(v).Float(), 'f', -1, 64)
	case bool:
		return strconv.FormatBool(d)
	default:
		b, _ := json.Marshal(v)
		return string(b)
	}
}

func MapStringToAny(req map[string]string) (res map[string]any) {
	res = make(map[string]any)
	for k, info := range req {
		res[k] = info
	}
	return res
}

func MapAnyToString(req map[string]any) (res map[string]string) {
	res = make(map[string]string)
	for k, info := range req {
		res[k] = ConvertAnyToStr(info)
	}
	return res
}
