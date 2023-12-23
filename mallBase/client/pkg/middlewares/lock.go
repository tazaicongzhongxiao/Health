package middlewares

import (
	"MyTestMall/mallBase/basics/pkg/app"
	"MyTestMall/mallBase/client/pkg/comm"
	"github.com/gin-gonic/gin"
)

// 并发锁
func VerifyLock(c *gin.Context) {
	if comm.GinRedLock(c) == false {
		comm.NewResponse(app.FailedToAcquireLock, nil, app.FailedToAcquireLockMessage).End(c)
		c.Abort()
	} else {
		c.Next()
	}
}
