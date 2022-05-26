package migrations

import (
	"errors"

	m "github.com/ferealqq/wienerlist/boardapi/models"
	s "github.com/ferealqq/wienerlist/seeders"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	// FIXME: not sure should we open & close the db connection here.
	// Is it a antipattern to be constantly closing db connections?
	return db.AutoMigrate(&m.Workspace{}, &m.Board{}, &m.Section{}, &m.Item{})
}

type TableToHandle struct {
	TableInterface interface{}
	SeederFunc     func(db *gorm.DB)
}

type TablesToHandle []TableToHandle

func (tables TablesToHandle) tableInterfaces() []interface{} {
	var list []interface{}
	for _, table := range tables {
		list = append(list, table.TableInterface)
	}
	return list
}

func MigrateSeedAfterwards(db *gorm.DB) {
	list := TablesToHandle{
		{
			SeederFunc:     s.SeedWorkspaces,
			TableInterface: &m.Workspace{},
		},
		{
			SeederFunc:     s.SeedBoards,
			TableInterface: &m.Board{},
		},
		{
			SeederFunc:     s.SeedSections,
			TableInterface: &m.Section{},
		},
		{
			SeederFunc:     s.SeedItems,
			TableInterface: &m.Item{},
		},
	}
	// FIXME should not panic, throw error instead
	if err := Migrate(db); err == nil && HasAllTables(db, list.tableInterfaces()...) {
		for _, table := range list {
			if err := db.First(table.TableInterface).Error; errors.Is(err, gorm.ErrRecordNotFound) {
				table.SeederFunc(db)
			}
		}
	} else {
		panic(err)
	}
}

func HasAllTables(db *gorm.DB, list ...interface{}) bool {
	// list of interfaces struct
	for _, table := range list {
		if !db.Migrator().HasTable(table) {
			return false
		}
	}
	return true
}
