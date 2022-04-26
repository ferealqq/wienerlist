package testing

import (
	"fmt"

	"gorm.io/gorm"
)

func ReinitTables(db *gorm.DB,tables []string, seedFuncs []func(*gorm.DB)) {
	for _, table := range tables {
		res := db.Exec("TRUNCATE TABLE "+table+" RESTART IDENTITY CASCADE;")
		fmt.Println(fmt.Sprintf(`Deleting table %s, affected rows: %v`,table,res.RowsAffected))
		if res.Error != nil {
			panic(res.Error)
		}
	}
	for _, seedFunc := range seedFuncs {
		seedFunc(db)
	}
}