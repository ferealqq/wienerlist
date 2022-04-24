package migrations

import (
	"errors"

	. "github.com/ferealqq/golang-trello-copy/server/boardapi/models"
	"github.com/ferealqq/golang-trello-copy/server/seeders"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	// FIXME: not sure should we open & close the db connection here.
	// Is it a antipattern to be constantly closing db connections?
	return db.AutoMigrate(&Board{})
}


func MigrateSeedAfterwards(db *gorm.DB){
	if err := Migrate(db); err == nil && db.Migrator().HasTable(&Board{}) {
		if err := db.First(&Board{}).Error; errors.Is(err, gorm.ErrRecordNotFound) {	
			seeders.SeedBoards(db)
		}
	}
}