package controller

import (
	"net/http"
	"strconv"

	"github.com/ferealqq/wienerlist/pkg/appenv"
	"github.com/ferealqq/wienerlist/pkg/database"
	"github.com/ferealqq/wienerlist/pkg/status"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type BaseController[M interface{}] struct {
	// Application environment
	AppEnv appenv.AppEnv
	// Connection to the database
	DB *gorm.DB
	// Context for the incoming request
	Context *gin.Context
}

func (b *BaseController[M]) SendJSON(status int, json interface{}) {
	b.Context.JSON(status, json)
}

func (b *BaseController[M]) SendInternalServerError(message string, err error) {
	log.WithFields(log.Fields{
		"env":    b.AppEnv.Env,
		"status": http.StatusInternalServerError,
		"error":  err,
	}).Error(message)
	b.SendJSON(http.StatusInternalServerError, gin.H{"message": message})
}

func (b *BaseController[M]) SendNotFound(message string) {
	b.SendJSON(http.StatusNotFound, status.Response{
		Status:  strconv.Itoa(http.StatusNotFound),
		Message: message,
	})
}

func (b *BaseController[M]) DefaultQueryInt(param string, def int) int {
	if val := b.Context.Query(param); val == "" {
		return def
	} else {
		if val, err := strconv.Atoi(val); err == nil {
			return val
		} else {
			return def
		}
	}
}

type UriId struct {
	ID uint `uri:"id" binding:"required,gt=0"`
}

func (b *BaseController[M]) GetUriId() (uint, error) {
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

// Get the post model of the request, incase of error send a bad request response
func (b *BaseController[M]) GetPostModel(m *M) error {
	if err := b.Context.ShouldBindJSON(&m); err != nil {
		// TODO Form a pattern in which we want to return error's
		b.SendJSON(http.StatusBadRequest, status.Response{
			Status:  strconv.Itoa(http.StatusBadRequest),
			Message: "malformed object",
		})
		return err
	}
	return nil
}

func MakeHandler[M interface{}](appEnv appenv.AppEnv, fn func(BaseController[M])) func(*gin.Context) {
	return func(c *gin.Context) {
		fn(BaseController[M]{
			DB:      database.DBConn,
			AppEnv:  appEnv,
			Context: c,
		})
	}
}
