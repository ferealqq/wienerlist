package boardapi

import (
	boardCtrl "github.com/ferealqq/golang-trello-copy/server/boardapi/controllers"
	app "github.com/ferealqq/golang-trello-copy/server/pkg/appenv"
	controllers "github.com/ferealqq/golang-trello-copy/server/pkg/controller"
	"github.com/gin-gonic/gin"
)

func BoardRouter(router *gin.Engine, appEnv app.AppEnv) {
	r := router.Group("/boards")
	{
		r.GET("/", controllers.MakeHandler(appEnv, boardCtrl.ListBoardsHandler))
		r.POST("/", controllers.MakeHandler(appEnv, boardCtrl.CreateBoardHandler))
		r.GET("/:id", controllers.MakeHandler(appEnv, boardCtrl.GetBoardHandler))
		r.PUT("/:id", controllers.MakeHandler(appEnv, boardCtrl.UpdateBoardHandler))
		r.DELETE("/:id", controllers.MakeHandler(appEnv, boardCtrl.DeleteBoardHandler))
	}
}
