package boardapi

import (
	"net/http"
	"strconv"

	"github.com/ferealqq/golang-trello-copy/server/migrations"
	"github.com/ferealqq/golang-trello-copy/server/pkg/database"
	"github.com/ferealqq/golang-trello-copy/server/pkg/status"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/unrolled/render"
	"gorm.io/gorm"
)

// FIXME: AppEnv renamed to controller?
type AppEnv struct {
	Render  *render.Render
	Version string
	Env     string
	Port    string
	// FIXME: Should we have db conn in the AppEnv?
	DBConn *gorm.DB
}

func (a *AppEnv) sendJSON(context *gin.Context, status int, json interface{}) {
	context.JSON(status, json)
}

func (a *AppEnv) sendInternalServerError(context *gin.Context, message string, err error) {
	log.WithFields(log.Fields{
		"env":    a.Env,
		"status": http.StatusInternalServerError,
		"error":  err,
	}).Error(message)
	a.sendJSON(context, http.StatusInternalServerError, gin.H{"message": message})
}

// FIXME move this to a pkg file
type UriId struct {
	ID uint `uri:"id" binding:"required,gt=0"`
}

func (a *AppEnv) getUriId(c *gin.Context) (uint, error) {
	var uri UriId
	if e := c.ShouldBindUri(&uri); e != nil {
		a.sendJSON(c, http.StatusBadRequest, status.Response{
			Status:  strconv.Itoa(http.StatusBadRequest),
			Message: "malformed id",
		})
		return 0, e
	}
	return uri.ID, nil
}

// CreateContextForTestSetup initialises an application context struct
// for testing purposes
func CreateContextForTestSetup() AppEnv {
	database.TestDBInit()
	// stop execution if migrations fail because tests won't be able to run
	if err := migrations.Migrate(database.DBConn); err != nil {
		panic(err)
	}
	testVersion := "0.0.0"
	// TODO init test database connection
	appEnv := AppEnv{
		Render:  render.New(),
		Version: testVersion,
		Env:     "LOCAL",
		Port:    "3001",
		DBConn:  database.DBConn,
	}
	return appEnv
}
