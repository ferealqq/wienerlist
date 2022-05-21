package testing

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"

	app "github.com/ferealqq/golang-trello-copy/server/pkg/appenv"
	ctrl "github.com/ferealqq/golang-trello-copy/server/pkg/controller"
	"github.com/ferealqq/golang-trello-copy/server/pkg/database"
	"github.com/gin-gonic/gin"
	"github.com/urfave/negroni"
	"gorm.io/gorm"
)

func ReinitTables(db *gorm.DB, tables []string, seedFuncs []func(*gorm.DB)) {
	for _, table := range tables {
		res := db.Exec("TRUNCATE TABLE " + table + " RESTART IDENTITY CASCADE;")
		fmt.Printf("\nDeleting table %s, affected rows: %v\n", table, res.RowsAffected)
		if res.Error != nil {
			panic(res.Error)
		}
	}
	for _, seedFunc := range seedFuncs {
		seedFunc(db)
	}
}

type HttpTestAction[M interface{}] struct {
	Name string
	// Functionality of the router
	RouterFunc func(*gin.Engine, app.AppEnv)
	// Path to which the request will be sent
	ReqPath string
	Handler func(ctrl.BaseController[M])
	Method  string
	Body    io.Reader
	Seeders []func(db *gorm.DB)
	Tables  []string
}

func (action *HttpTestAction[M]) Run() *httptest.ResponseRecorder {
	appenv := app.CreateTestAppEnv()
	//FIXME figure out a better way to give out the table name, table names could change so this is a little problematic approach
	ReinitTables(database.DBConn, action.Tables, action.Seeders)
	r, _ := http.NewRequest(action.Method, action.ReqPath, action.Body)
	w := httptest.NewRecorder()
	router := gin.Default()
	// recreate route from routes.go
	action.RouterFunc(router, appenv)
	n := negroni.New()
	n.UseHandler(router)
	n.ServeHTTP(w, r)

	return w
}
