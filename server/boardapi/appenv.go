package boardapi

import (
	"github.com/ferealqq/golang-trello-copy/server/migrations"
	"github.com/ferealqq/golang-trello-copy/server/pkg/database"
	"github.com/unrolled/render"
	"gorm.io/gorm"
)

// AppEnv holds application configuration data
type AppEnv struct {
	Render      *render.Render
	Version     string
	Env         string
	Port        string
	// FIXME: Should we have db conn in the AppEnv? 
	DBConn	    *gorm.DB
}

// CreateContextForTestSetup initialises an application context struct
// for testing purposes
func CreateContextForTestSetup() AppEnv {
	database.TestDBInit()
	migrations.Migrate(database.DBConn)
	testVersion := "0.0.0"
	// TODO init test database connection
	appEnv := AppEnv{
		Render:    render.New(),
		Version:   testVersion,
		Env:       "LOCAL",
		Port:      "3001",
		DBConn:    database.DBConn,
	}
	return appEnv
}
