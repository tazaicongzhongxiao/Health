package base

import (
	"clientHealthy/controllers"
	"github.com/gin-gonic/gin"
)

func InitBaseRouter(Router *gin.RouterGroup) gin.IRoutes {
	router := Router.Group("base")
	{
		router.GET("bodyParam/page", controllers.BodyParam)
		router.GET("bodyParam/list", controllers.BodyParamList)
		router.POST("bodyParam/save", controllers.BodyParamSave)
		router.DELETE("bodyParam/delete", controllers.BodyParamDelete)
		router.GET("food/page", controllers.Food)
		router.GET("food/list", controllers.FoodList)
		router.POST("food/save", controllers.FoodSave)
		router.DELETE("food/delete", controllers.FoodDelete)
		router.GET("sports/page", controllers.Sports) //身体参数分页
		router.GET("sports/list", controllers.SportsList)
		router.POST("sports/save", controllers.SportsSave)
		router.DELETE("sports/delete", controllers.SportsDelete)
	}
	return router
}
