package pic

import (
	"github.com/gin-gonic/gin"
	"gitlab.mall.com/mallBase/basics/pkg/app"
)

func InitPicRouter(engine *gin.Engine) {
	engine.GET("Pic/:name/:id", func(c *gin.Context) {
		name := c.Param("name")
		id := c.Param("id")
		res := "https://127.0.0.01:31000" + "/Pic" + name + id + ".png"
		app.Println(3, res)
		c.File(res)
	})
}
