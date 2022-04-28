package container

import (
	"net/http"
	"strconv"

	"github.com/ferealqq/golang-trello-copy/server/migrations"
	"github.com/ferealqq/golang-trello-copy/server/pkg/database"
	"github.com/ferealqq/golang-trello-copy/server/pkg/status"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// FIXME: AppEnv renamed to controller?
type AppContainer struct {
	Version string
	Env     string
	Port    string
	// FIXME: Should we have db conn in the AppEnv?
	DBConn *gorm.DB
}

func (a *AppContainer) SendJSON(context *gin.Context, status int, json interface{}) {
	context.JSON(status, json)
}

func (a *AppContainer) SendInternalServerError(context *gin.Context, message string, err error) {
	log.WithFields(log.Fields{
		"env":    a.Env,
		"status": http.StatusInternalServerError,
		"error":  err,
	}).Error(message)
	a.SendJSON(context, http.StatusInternalServerError, gin.H{"message": message})
}

// FIXME move this to a pkg file
type UriId struct {
	ID uint `uri:"id" binding:"required,gt=0"`
}

func (a *AppContainer) GetUriId(c *gin.Context) (uint, error) {
	var uri UriId
	if e := c.ShouldBindUri(&uri); e != nil {
		a.SendJSON(c, http.StatusBadRequest, status.Response{
			Status:  strconv.Itoa(http.StatusBadRequest),
			Message: "malformed id",
		})
		return 0, e
	}
	return uri.ID, nil
}

// CreateContextForTestSetup initialises an application context struct
// for testing purposes
func CreateContextForTestSetup() AppContainer {
	database.TestDBInit()
	// stop execution if migrations fail because tests won't be able to run
	if err := migrations.Migrate(database.DBConn); err != nil {
		panic(err)
	}
	testVersion := "0.0.0"
	// TODO init test database connection
	appContainer := AppContainer{
		Version: testVersion,
		Env:     "LOCAL",
		Port:    "3001",
		DBConn:  database.DBConn,
	}
	return appContainer
}

func MakeHandler(appContainer AppContainer, fn func(*gin.Context, AppContainer)) func(*gin.Context) {
	return func(c *gin.Context) {
		fn(c, appContainer)
	}
}
