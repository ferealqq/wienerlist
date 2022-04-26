package migrations

import (
	"errors"

	. "github.com/ferealqq/golang-trello-copy/server/boardapi/models"
	. "github.com/ferealqq/golang-trello-copy/server/seeders"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	// FIXME: not sure should we open & close the db connection here.
	// Is it a antipattern to be constantly closing db connections?
	return db.AutoMigrate(&Board{}, &Section{})
}
type TableToHandle struct {
	TableInterface interface{}
	SeederFunc func(db *gorm.DB)
}

type TablesToHandle []TableToHandle

func (tables TablesToHandle) tableInterfaces() []interface{} {
	var list []interface{}
	for _, table := range tables {
		list = append(list, table.TableInterface)
	}
	return list
}

func MigrateSeedAfterwards(db *gorm.DB){
	list := TablesToHandle{
		{
			SeederFunc: SeedBoards,
			TableInterface: &Board{},
		},
		{
			SeederFunc: SeedSections,
			TableInterface: &Section{},
		},
	}
	if err := Migrate(db); err == nil && HasAllTables(db, list.tableInterfaces()...) {
		for _, table := range list {
			if err := db.First(table.TableInterface).Error; errors.Is(err, gorm.ErrRecordNotFound) {	
				table.SeederFunc(db)
			}
		}
	}
}

func HasAllTables(db *gorm.DB,list ...interface{}) bool {
	// list of interfaces struct
	for _, table := range list {
		db.Migrator()
		if !db.Migrator().HasTable(table) {
			return false
		}
	}
	return true
}