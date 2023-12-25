package response

import (
	"MyTestMall/mallBase/basics/pkg/app"
	"fmt"
	json "github.com/json-iterator/go"
)

func PicResponse(data interface{}) interface{} {
	var result struct {
		Code    int32  `json:"code"`    // 状态码,这个状态码是与前端和APP约定的状态码,非HTTP状态码
		Data    []byte `json:"data"`    // 返回数据
		Message string `json:"message"` // 自定义返回的消息内容
	}
	_ = app.Unmarshal(data, &result)
	if result.Code == 0 {
		result.Code = app.Success
	}
	if result.Message == "" {
		result.Message = app.SuccessMessage
	}
	if len(result.Data) != 0 {
		var info ResFood
		_ = json.Unmarshal(result.Data, &info)
		for _, v := range info.Pic {
			v = fmt.Sprintf("http://192.168.2.139:31000/Pic/%s/%s.png", info.Name, v)
		}
		_ = app.Unmarshal(info, &result)
	}
	return result
}
