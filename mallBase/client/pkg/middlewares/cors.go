package middlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// CORS 全局允许跨域
func CORS(c *gin.Context) {
	//if c.Request.Header.Get(`X-Requested-With`) != "" || c.Request.Header.Get(`Origin`) != "" {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Access-Control-Max-Age", "86400")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,UPDATE")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "*")
	// c.Writer.Header().Set("Access-Control-Allow-Headers", "Access-Control-Allow-Origin,User-Agent,Referer,Accept,Sec-Fetch-Mode,Origin,Content-Type,Content-Length,Accept-Encoding,x-requested-with,StoreId,Platform,Authorization")
	c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	//}
	if c.Request.Method == "OPTIONS" {
		c.AbortWithStatus(http.StatusNoContent)
		c.Abort()
		return
	}
	c.Next()
}
