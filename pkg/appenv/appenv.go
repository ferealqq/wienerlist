package appenv

import (
	"github.com/ferealqq/golang-trello-copy/server/migrations"
	"github.com/ferealqq/golang-trello-copy/server/pkg/database"
)

type AppEnv struct {
	Version string
	Env     string
	Port    string
}

func CreateTestAppEnv() AppEnv {
	database.TestDBInit()
	// stop execution if migrations fail because tests won't be able to run
	if err := migrations.Migrate(database.DBConn); err != nil {
		panic(err)
	}

	testVersion := "0.0.0"

	return AppEnv{
		Version: testVersion,
		Env:     "LOCAL",
		Port:    "3001",
	}
}
