package boardapi

import (
	boardCtrl "github.com/ferealqq/golang-trello-copy/server/boardapi/controllers"
	"github.com/ferealqq/golang-trello-copy/server/boardapi/models"
	app "github.com/ferealqq/golang-trello-copy/server/pkg/appenv"
	controllers "github.com/ferealqq/golang-trello-copy/server/pkg/controller"
	"github.com/gin-gonic/gin"
)

func BoardRouter(router *gin.Engine, appEnv app.AppEnv) {
	r := router.Group("/boards")
	{
		r.GET("/", controllers.MakeHandler[models.Board](appEnv, boardCtrl.ListBoardsHandler))
		r.POST("/", controllers.MakeHandler[models.Board](appEnv, boardCtrl.CreateBoardHandler))
		r.GET("/:id", controllers.MakeHandler[models.Board](appEnv, boardCtrl.GetBoardHandler))
		r.PUT("/:id", controllers.MakeHandler[models.Board](appEnv, boardCtrl.UpdateBoardHandler))
		r.DELETE("/:id", controllers.MakeHandler[models.Board](appEnv, boardCtrl.DeleteBoardHandler))
	}
}
