package boardapi

import "github.com/gin-gonic/gin"

func BoardRouter(router *gin.Engine, appEnv AppEnv) {
	r := router.Group("/boards")
	{
		r.GET("/", MakeHandler(appEnv, ListBoardsHandler))
		r.POST("/", MakeHandler(appEnv, CreateBoardHandler))
		r.GET("/:id", MakeHandler(appEnv, GetBoardHandler))
		r.PUT("/:id", MakeHandler(appEnv, UpdateBoardHandler))
		r.DELETE("/:id", MakeHandler(appEnv, DeleteBoardHandler))
	}
}
