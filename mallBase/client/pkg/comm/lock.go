package comm

import (
	"MyTestMall/mallBase/basics/pkg/dredis"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// GIN 并发处理检查
func GinRedLock(c *gin.Context) bool {
	if c.Request.Method == http.MethodPut || c.Request.Method == http.MethodDelete {
		name := "redLock:"
		if Uid, ok := c.Get("uid"); ok {
			name += strconv.Itoa(int(Uid.(int64)))
		} else {
			name += c.ClientIP() + ":" + c.Request.URL.Path
		}
		lock := dredis.NewRedisLockWithParam(name, 5, 2, 300)
		if err := lock.Lock(); err != nil {
			return false
		}
		c.Set("lock", lock)
	}
	return true
}

// GIN 解除并发锁
func GinUnLock(c *gin.Context) bool {
	if value, exists := c.Get("lock"); exists {
		if lock, ok := value.(*dredis.RedisLock); ok {
			_ = lock.UnLock()
		}
	}
	return true
}
