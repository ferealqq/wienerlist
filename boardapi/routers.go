package boardapi

import (
	ctrl "github.com/ferealqq/wienerlist/boardapi/controllers"
	app "github.com/ferealqq/wienerlist/pkg/appenv"
	controllers "github.com/ferealqq/wienerlist/pkg/controller"
	"github.com/gin-gonic/gin"
)

func WorkspaceRouter(router *gin.RouterGroup, appEnv app.AppEnv) {
	r := router.Group("/workspaces")
	{
		r.GET("/", controllers.MakeHandler(appEnv, ctrl.ListWorkspacesHandler))
		r.POST("/", controllers.MakeHandler(appEnv, ctrl.CreateWorkspaceHandler))
		r.GET("/:id", controllers.MakeHandler(appEnv, ctrl.GetWorkspaceHandler))
		r.PATCH("/:id", controllers.MakeHandler(appEnv, ctrl.UpdateWorkspaceHandler))
		r.DELETE("/:id", controllers.MakeHandler(appEnv, ctrl.DeleteWorkspaceHandler))
	}
}

func BoardRouter(router *gin.RouterGroup, appEnv app.AppEnv) {
	r := router.Group("/boards")
	{
		r.GET("/", controllers.MakeHandler(appEnv, ctrl.ListBoardsHandler))
		r.POST("/", controllers.MakeHandler(appEnv, ctrl.CreateBoardHandler))
		r.GET("/:id", controllers.MakeHandler(appEnv, ctrl.GetBoardHandler))
		r.PATCH("/:id", controllers.MakeHandler(appEnv, ctrl.UpdateBoardHandler))
		r.DELETE("/:id", controllers.MakeHandler(appEnv, ctrl.DeleteBoardHandler))
	}
}

func SectionRouter(router *gin.RouterGroup, appEnv app.AppEnv) {
	r := router.Group("/sections")
	{
		r.GET("/", controllers.MakeHandler(appEnv, ctrl.ListSectionsHandler))
		r.POST("/", controllers.MakeHandler(appEnv, ctrl.CreateSectionHandler))
		r.GET("/:id", controllers.MakeHandler(appEnv, ctrl.GetSectionHandler))
		r.PATCH("/:id", controllers.MakeHandler(appEnv, ctrl.UpdateSectionHandler))
		r.DELETE("/:id", controllers.MakeHandler(appEnv, ctrl.DeleteSectionHandler))
	}
}

func ItemRouter(router *gin.RouterGroup, appEnv app.AppEnv) {
	r := router.Group("/items")
	{
		r.GET("/", controllers.MakeHandler(appEnv, ctrl.ListItemsHandler))
		r.POST("/", controllers.MakeHandler(appEnv, ctrl.CreateItemHandler))
		r.PATCH("/:id", controllers.MakeHandler(appEnv, ctrl.UpdateItemHandler))
		r.GET("/:id", controllers.MakeHandler(appEnv, ctrl.GetItemHandler))
		r.DELETE("/:id", controllers.MakeHandler(appEnv, ctrl.DeleteItemHandler))
	}
}
