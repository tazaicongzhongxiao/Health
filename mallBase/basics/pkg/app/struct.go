package app

import (
	"MyTestMall/mallBase/basics/tools/contains"
	"context"
	"github.com/TarsCloud/TarsGo/tars/util/current"
	json "github.com/json-iterator/go"
	"reflect"
	"strconv"
	"strings"
)

func SetContext(Name string, value interface{}, filter ...string) (res map[string]string) {
	res = make(map[string]string)
	if Name == "elect" {
		// 传递查询的数据库字段
		res[Name] = GetStructName(value, filter)
	} else {
		if v, ok := value.(string); ok {
			res[Name] = v
		}
	}
	return
}

func GetContext(ctx context.Context, Name string) string {
	req, _ := current.GetRequestContext(ctx)
	if req[Name] != "" {
		return req[Name]
	}
	return ""
}

// GetStructName
// 获取 proto 结构体 JSON 名称
// filter 读取结构体 JSON 后排除的字段名
func GetStructName(structName interface{}, filter []string) string {
	t := reflect.TypeOf(structName)
	result := make([]string, 0)
	if t.Kind() == reflect.String {
		return structName.(string)
	} else if t.Kind() == reflect.Ptr {
		var t interface{}
		_ = json.Unmarshal(Struct2Json(structName), &t)
		if len(filter) > 0 {
			for k, _ := range t.(map[string]interface{}) {
				if contains.ContainsString(filter, k) == -1 {
					result = append(result, k)
				}
			}
		} else {
			for k, _ := range t.(map[string]interface{}) {
				result = append(result, k)
			}
		}
	} else if t.Kind() == reflect.Slice {
		s := reflect.ValueOf(structName)
		if len(filter) > 0 {
			for i := 0; i < s.Len(); i++ {
				jsonName := s.Index(i).Interface().(string)
				if contains.ContainsString(filter, jsonName) == -1 {
					result = append(result, jsonName)
				}
			}
		} else {
			for i := 0; i < s.Len(); i++ {
				result = append(result, s.Index(i).Interface().(string))
			}
		}
	} else if t.Kind() == reflect.Struct {
		fieldNum := t.NumField()
		if len(filter) > 0 {
			for i := 0; i < fieldNum; i++ {
				jsonName := t.Field(i).Tag.Get("json")
				if jsonName != "-" && jsonName != "" {
					comma := strings.Index(jsonName, ",")
					if comma > 1 {
						jsonName = jsonName[:comma]
					}
					if contains.ContainsString(filter, jsonName) == -1 {
						result = append(result, jsonName)
					}
				}
			}
		} else {
			for i := 0; i < fieldNum; i++ {
				jsonName := t.Field(i).Tag.Get("json")
				if jsonName != "-" && jsonName != "" {
					comma := strings.Index(jsonName, ",")
					if comma > 1 {
						jsonName = jsonName[:comma]
					}
					result = append(result, jsonName)
				}
			}
		}
	}
	return strings.Join(result, ", ")
}

func Struct2Json(form interface{}) []byte {
	jsons, err := json.Marshal(form)
	if err != nil {
		jsons, err = json.Marshal(map[string]string{"err": err.Error()})
		if err != nil {
			return []byte(err.Error())
		}
	}
	return jsons
}

func UnmarshalJson(req interface{}, v ...interface{}) (err error) {
	jsonb := Struct2Json(req)
	for _, val := range v {
		if err = json.Unmarshal(jsonb, &val); err != nil {
			return err
		}
	}
	return nil
}

// Unmarshal
// @Description: 读取值转换
func Unmarshal(req interface{}, v interface{}) error {
	return json.Unmarshal(Struct2Json(req), &v)
}

// UnmarshalElect
// @Description: 读取值转换
func UnmarshalElect(req interface{}, v interface{}) (string, error) {
	return GetStructName(req, nil), json.Unmarshal(Struct2Json(req), &v)
}

// GetBool
// @Description: 根据 interface 获取值
func GetBool(str interface{}, empty ...bool) bool {
	switch str.(type) {
	case string:
		return str.(string) == "true"
	case bool:
		return str.(bool)
	default:
		if empty == nil {
			return false
		}
		return empty[0]
	}
}

func GetString(str interface{}, empty ...string) string {
	switch str.(type) {
	case string:
		return str.(string)
	case float64:
		return strconv.Itoa(int(str.(float64)))
	default:
		if empty == nil {
			return ""
		}
		return empty[0]
	}
}

func GetStrings(str interface{}, empty ...string) (res []string) {
	if err := Unmarshal(str, &res); err != nil {
		if empty == nil {
			return nil
		}
		return empty
	}
	return
}

func GetFloat(str interface{}, empty ...float64) float64 {
	switch str.(type) {
	case string:
		s, _ := strconv.ParseFloat(str.(string), 64)
		return s
	case float64:
		return str.(float64)
	default:
		if empty == nil {
			return 0
		}
		return empty[0]
	}
}

func GetFloats(str interface{}, empty ...float64) (value []float64) {
	if val, ok := str.(interface{}); ok {
		_ = Unmarshal(val, &value)
		return value
	} else {
		if empty == nil {
			return nil
		}
		return empty
	}
}
