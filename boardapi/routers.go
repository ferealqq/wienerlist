package boardapi

import (
	. "github.com/ferealqq/golang-trello-copy/server/boardapi/controllers"
	. "github.com/ferealqq/golang-trello-copy/server/pkg/container"
	"github.com/gin-gonic/gin"
)

func BoardRouter(router *gin.Engine, appC AppContainer) {
	r := router.Group("/boards")
	{
		r.GET("/", MakeHandler(appC, ListBoardsHandler))
		r.POST("/", MakeHandler(appC, CreateBoardHandler))
		r.GET("/:id", MakeHandler(appC, GetBoardHandler))
		r.PUT("/:id", MakeHandler(appC, UpdateBoardHandler))
		r.DELETE("/:id", MakeHandler(appC, DeleteBoardHandler))
	}
}
