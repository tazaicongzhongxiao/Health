package middlewares

import (
	"MyTestMall/mallBase/basics/pkg/app"
	"MyTestMall/mallBase/basics/tools/dstring"
	"MyTestMall/mallBase/client/pkg/comm"
	"MyTestMall/mallBase/client/pkg/jwt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

var (
	// JwtHeaderKey jwt token 在HTTP请求中的header名
	JwtHeaderKey = "Authorization"
)

func GetAuth(c *gin.Context) (uid int64, Type int8) {
	if token := c.GetHeader(JwtHeaderKey); token != "" {
		claims, err := jwt.ParseToken(token)
		if err == nil {
			return claims.Uid, claims.Type
		}
	}
	return
}

// VerifyAuth
// @Description: 中间件 验证用户有效性的中间件
func VerifyAuth(c *gin.Context) {
	token := c.GetHeader(JwtHeaderKey)
	if token == "" {
		comm.NewResponse(app.AuthFail, nil, app.AuthFailMessage).End(c, http.StatusUnauthorized)
	} else {
		claims, err := jwt.ParseToken(token)
		if err != nil {
			comm.NewResponse(app.AuthFail, nil, err.Error()).End(c, http.StatusUnauthorized)
		} else {
			c.Set("uid", claims.Uid)
			c.Set("uType", claims.Type)
			c.Header(JwtHeaderKey, token)
			c.Next()
			return
		}
	}
	c.Abort()
}

// GetVerifyAuth
// @Description: 验证用户有效性的中间件
// @param c
func GetVerifyAuth(c *gin.Context) {
	if token := c.GetHeader(JwtHeaderKey); token != "" {
		claims, err := jwt.ParseToken(token)
		if err == nil {
			c.Set("uid", claims.Uid)
			c.Set("uType", claims.Type)
			c.Header(JwtHeaderKey, token)
			c.Next()
			return
		}
	}
}

// GetAgentPlatform
// @Description: 通过浏览器类型获取平台标识ID
func GetAgentPlatform(c *gin.Context) int32 {
	agent := c.GetHeader("User-Agent")
	if dstring.MatchString(`(?i:MicroMessenger)`, agent) {
		return 5
	} else if dstring.MatchString(`(?i:Alipay)`, agent) {
		return 6
	} else if dstring.MatchString(`(?i:Mobile|iPod|iPhone|Android|Opera Mini|BlackBerry|webOS|UCWEB|Blazer|PSP)`, agent) {
		return 2
	} else {
		return 1
	}
}

func GetPlatform(platform string) int32 {
	switch platform {
	case "web":
		return 1
	case "h5":
		return 2
	case "rn": // APP
		return 3
	case "wechat": // 微信公众号
		return 4
	case "weapp": // 微信小程序
		return 5
	case "alipay": // 支付宝
		return 6
	case "dy": // 抖音小程序
		return 7
	case "quickapp": // 快手小程序
		return 8
	case "qq": // QQ小程序
		return 9
	}
	return 0
}

// GetUid
// @Description: 获取当前JWT 用户ID
// @param c
// @return id
func GetUid(c *gin.Context) (id int64) {
	if tmp, ok := c.Get("uid"); ok {
		return tmp.(int64)
	}
	return 0
}

// GetType
// @Description: 获取当前JWT 用户类型
// @param c
// @return Type
func GetType(c *gin.Context) (Type int8) {
	if tmp, ok := c.Get("uType"); ok {
		return tmp.(int8)
	}
	return 0
}

// GetCoId
// @Description: 获取当前JWT 企业ID
// @param c
// @return coId
func GetCoId(c *gin.Context) (coId int64) {
	if tmp, ok := c.Get("coId"); ok {
		return tmp.(int64)
	}
	return 0
}

// GetStoreId
// @Description: 获取当前JWT 店铺ID
// @param c
// @return storeId
func GetStoreId(c *gin.Context) (storeId int64) {
	storeId, _ = strconv.ParseInt(c.Param("store_id"), 10, 64)
	if storeId == 0 {
		if jwfStoreId := c.GetHeader("storeid"); jwfStoreId != "" {
			storeId, _ = strconv.ParseInt(jwfStoreId, 10, 64)
		}
	}
	return storeId
}

func GetFastSource(c *gin.Context) (storeId, uid int64) {
	if jwfStoreId := c.GetHeader("storeid"); jwfStoreId != "" {
		storeId, _ = strconv.ParseInt(jwfStoreId, 10, 64)
	} else {
		storeId, _ = strconv.ParseInt(c.Param("store_id"), 10, 64)
	}
	if jwfUid, ok := c.Get("uid"); ok {
		uid = jwfUid.(int64)
	}
	return storeId, uid
}
