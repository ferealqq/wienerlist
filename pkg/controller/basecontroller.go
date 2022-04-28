package controller

import (
	"net/http"
	"strconv"

	"github.com/ferealqq/golang-trello-copy/server/boardapi/models"
	"github.com/ferealqq/golang-trello-copy/server/pkg/appenv"
	"github.com/ferealqq/golang-trello-copy/server/pkg/database"
	"github.com/ferealqq/golang-trello-copy/server/pkg/status"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// FIXME: BaseController renamed to ControllerContainer/BaseController ? BaseController could be misleading.
type BaseController struct {
	// Application environment
	AppEnv appenv.AppEnv
	// Connection to the database
	DB *gorm.DB
	// Context for the incoming request
	Context *gin.Context
}

func (a *BaseController) SendJSON(status int, json interface{}) {
	a.Context.JSON(status, json)
}

func (a *BaseController) SendInternalServerError(message string, err error) {
	log.WithFields(log.Fields{
		"env":    a.AppEnv.Env,
		"status": http.StatusInternalServerError,
		"error":  err,
	}).Error(message)
	a.SendJSON(http.StatusInternalServerError, gin.H{"message": message})
}

// FIXME move this to a pkg file
type UriId struct {
	ID uint `uri:"id" binding:"required,gt=0"`
}

func (b *BaseController) GetUriId() (uint, error) {
	var uri UriId
	if e := b.Context.ShouldBindUri(&uri); e != nil {
		b.SendJSON(http.StatusBadRequest, status.Response{
			Status:  strconv.Itoa(http.StatusBadRequest),
			Message: "malformed id",
		})
		return 0, e
	}
	return uri.ID, nil
}

// FIXME: Use generics with this function instead of hardcoded models.Board when project has been upgraded to go1.18
func (b *BaseController) GetPostModel(m models.Board) (models.Board, error) {
	if err := b.Context.ShouldBindJSON(&m); err != nil {
		// TODO Form a pattern in which we want to return error's
		b.SendJSON(http.StatusBadRequest, status.Response{
			Status:  strconv.Itoa(http.StatusBadRequest),
			Message: "malformed object",
		})
		return m, err
	}
	return m, nil
}

func MakeHandler(appEnv appenv.AppEnv, fn func(BaseController)) func(*gin.Context) {
	return func(c *gin.Context) {
		fn(BaseController{
			DB:      database.DBConn,
			AppEnv:  appEnv,
			Context: c,
		})
	}
}
