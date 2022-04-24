package migrations

import (
	models "github.com/ferealqq/golang-trello-copy/server/boardapi/models"
	database "github.com/ferealqq/golang-trello-copy/server/pkg/database"
)

func Migrate(){
	// FIXME: not sure should we open & close the db connection here.
	// Is it a antipattern to be constantly closing db connections?
	database.DBConn.AutoMigrate(&models.Board{})
}