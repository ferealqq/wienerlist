package testing

import (
	"fmt"

	"gorm.io/gorm"
)

func ReinitTable(db *gorm.DB,table string, seedFunc func(*gorm.DB)) {
	res := db.Exec("TRUNCATE TABLE "+table+" RESTART IDENTITY CASCADE;")
	fmt.Println(fmt.Sprintf(`Deleting table %s, affected rows: %v`,table,res.RowsAffected))
	if res.Error != nil {
		panic(res.Error)
	}
	seedFunc(db)
}