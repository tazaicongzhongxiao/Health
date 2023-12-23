package comm

import (
	"MyTestMall/mallBase/basics/pkg/app"
	"github.com/TarsCloud/TarsGo/tars"
	"github.com/gin-gonic/gin"
	json "github.com/json-iterator/go"
	"net/http"
)

// Response HTTP返回数据结构体, 可使用这个, 也可以自定义
type Response struct {
	Code    int32       `json:"code"`    // 状态码,这个状态码是与前端和APP约定的状态码,非HTTP状态码
	Data    interface{} `json:"data"`    // 返回数据
	Message string      `json:"message"` // 自定义返回的消息内容
}

// ResponseTmlMsg Response 模板输出内容
type ResponseTmlMsg struct {
	Name string `json:"name"` // 标题
	Icon string `json:"icon"` // 主体、图标 success warning info error
	Time int8   `json:"time"` // 自动关闭时间
	Jump string `json:"jump"` // 跳转地址
	Err  string `json:"err"`  // 错误提示
}

func init() {
	// INFO < DEBUG < WARN < ERROR < NONE
	if app.Cfg.LogLevel == "NONE" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}
}

// End 在调用了这个方法之后,还是需要 return 的
func (rsp *Response) End(c *gin.Context, httpStatus ...int) {
	status := http.StatusOK
	if len(httpStatus) > 0 {
		status = httpStatus[0]
	}
	// 如果存在并发锁 输出时解锁
	GinUnLock(c)
	c.JSON(status, rsp)
}

// NewResponse 接口返回统一使用这个
//
//	code 服务端与客户端和web端约定的自定义状态码
//	data 具体的返回数据
//	message 可不传,自定义消息
func NewResponse(code int32, data interface{}, message ...string) *Response {
	msg := ""
	if len(message) > 0 {
		msg = message[0]
	}
	// fmt.Println(i18n.Get().NewPrinter(language.Chinese).Translate("Hello", nil, "one"))
	return &Response{Code: code, Data: data, Message: msg}
}

// ValidatorResponse
// @Description: 输出表单验证失败提示
func ValidatorResponse(ctx *gin.Context, err interface{}) {
	NewResponse(app.Validator, err, app.ValidatorMessage).End(ctx)
}

func ApiResponseByte(body []byte) (v interface{}) {
	_ = json.Unmarshal(body, &v)
	return
}

// ApiResponse
// @Description:  信息输出
// tips 提示类型 0 仅输出 1 消息提示 2 通知 3 弹框 tips 大于100 为输出自定义 code 编号
// errTips 如果无错误是否忽略弹出提示
// ignore 如果为弹出窗正确输出 JSON数据 错误输出弹窗
func ApiResponse(ctx *gin.Context, data interface{}, err error) {
	if err == nil {
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
			data = ApiResponseByte(result.Data)
		}
		NewResponse(result.Code, data, result.Message).End(ctx)
	} else {
		var result Response
		result.Code = tars.GetErrorCode(err)
		switch result.Code {
		case app.Fail, 1:
			result.Message = err.Error()
			switch result.Message {
			case "mongo: no documents in result", "record not found":
				result.Code = app.NotExist
				result.Message = app.NotExistMessage
			default:
				result.Code = 1
				if gin.Mode() == gin.ReleaseMode {
					result.Code = app.RequestError
					result.Message = app.RequestMessage
				}
			}
		case app.Success:
			_ = json.Unmarshal([]byte(err.Error()), &data)
		case app.Validator:
			result.Message = app.ValidatorMessage
			_ = json.Unmarshal([]byte(err.Error()), &data)
		default:
			if gin.Mode() == gin.ReleaseMode {
				result.Code = app.RequestError
				result.Message = app.RequestMessage
			} else {
				result.Message = err.Error()
			}
		}
		NewResponse(result.Code, data, result.Message).End(ctx)
	}
	return
}

func TplResponse(ctx *gin.Context, name string, err error) {
	ctx.HTML(http.StatusOK, name, ResponseTmlMsg{
		Err: err.Error(),
	})
}
