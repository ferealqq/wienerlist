package boardapi

import (
	"net/http"
	"strconv"

	"github.com/ferealqq/golang-trello-copy/server/migrations"
	"github.com/ferealqq/golang-trello-copy/server/pkg/database"
	"github.com/ferealqq/golang-trello-copy/server/pkg/status"
	log "github.com/sirupsen/logrus"
	"github.com/unrolled/render"
	"gorm.io/gorm"
)

// AppEnv holds application configuration data
type AppEnv struct {
	Render  *render.Render
	Version string
	Env     string
	Port    string
	// FIXME: Should we have db conn in the AppEnv?
	DBConn *gorm.DB
}

// Send JSON response to client if able to. In case of error return ISE
func (a *AppEnv) sendJSON(w http.ResponseWriter, status int, json interface{}) {
	if err := a.Render.JSON(w, status, json); err != nil {
		a.sendInternalServerError(w, "Error sending JSON response", err)
	}
}

// Send ISE response to client in case of error log it
func (a *AppEnv) sendInternalServerError(w http.ResponseWriter, message string, err error) {
	response := status.Response{
		Status:  strconv.Itoa(http.StatusInternalServerError),
		Message: message,
	}

	log.WithFields(log.Fields{
		"env":    a.Env,
		"status": http.StatusInternalServerError,
		"error":  err,
	}).Error(message)

	// If the AppEnv.Render.JSON() call fails, we need to send a 500 response and just log the error. We can't do more at this point
	if newErr := a.Render.JSON(w, http.StatusInternalServerError, response); newErr != nil {
		log.WithFields(log.Fields{
			"env":   a.Env,
			"error": newErr,
		}).Error("Couldn't send ISE response to client. Render could be broken")
	}
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
