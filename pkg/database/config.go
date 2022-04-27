package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DBConn *gorm.DB
)

/**
 * Initialize the database connection
 */
func InitDB() {
	var err error
	// TODO check that the database connection is closed or not initialized
	dsn := "host=localhost user=postgres password=postgres dbname=BOARD port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	DBConn, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(err)
	}
}

/**
* Close the current database connection
 */
func Close() {
	db, err := DBConn.DB()

	if err != nil {
		panic(err)
	}

	db.Close()
}

func TestDBInit() {
	var err error
	// TODO check that the database connection is closed or not initialized
	// FIXME Database configurations to env
	dsn := "host=localhost user=postgres password=postgres dbname=TEST_BOARD port=5433 sslmode=disable TimeZone=Asia/Shanghai"
	DBConn, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(err)
	}
}
