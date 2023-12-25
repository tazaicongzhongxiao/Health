package main

import (
	"MyTestMall/mallBase/basics/pkg/app"
	"MyTestMall/mallBase/basics/pkg/dredis"
	"MyTestMall/mallBase/client/pkg/jwt"
	"MyTestMall/mallBase/client/pkg/middlewares"
	"MyTestMall/mallBase/client/pkg/validator"
	"clientHealthy/router/base"
	"fmt"
	"github.com/TarsCloud/TarsGo/tars"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"net/http"
)

func init() {
	jwt.Start()
	dredis.Start()
	binding.Validator = new(validator.Validator)
}

func main() {
	mux := &tars.TarsHttpMux{}
	engine := mux.GetGinEngine()
	engine.Use(middlewares.CORS)
	engine.StaticFS("/Pic", http.Dir("./pic"))
	engine.GET("Pic", func(c *gin.Context) {
		name := c.Query("name")
		id := c.Query("id")
		res := fmt.Sprintf("http://192.168.2.139:31000/Pic/%s/%s.png", name, id)
		app.Println(3, res)
		c.JSON(http.StatusOK, res)
	})
	ApiGroup := engine.Group("/api/v1/healthy")
	if app.Cfg.LogLevel == "DEBUG" || app.Cfg.LogLevel == "INFO" {
		gin.SetMode(gin.DebugMode)
		ApiGroup.GET("/doc/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
		// 开启链路
		//cf := zipkintracing.ZipkinClientFilter()
		//tars.RegisterClientFilter(cf)
		//zipkintracing.InitV2()
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	ApiGroup.Use(middlewares.VerifyLock)
	base.InitBaseRouter(ApiGroup) // 用户路由
	tars.AddHttpServant(mux, app.Cfg.App+"."+app.Cfg.Server+".healthyObj")
	tars.Run()
}
