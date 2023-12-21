package app

import json "github.com/json-iterator/go"

func Struct2Json(req interface{}) []byte {
	result, err := json.Marshal(req)
	if err != nil {
		result, err = json.Marshal(map[string]string{"err": err.Error()})
		if err != nil {
			return []byte(err.Error())
		}
	}
	return result
}
