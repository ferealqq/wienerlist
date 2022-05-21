package database

import (
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DBConn *gorm.DB
)

/**
 * Initialize the database connection
 */
func InitDB() error {
	var err error
	// TODO check that the database connection is closed or not initialized
	dsn := os.Getenv("DB_DSN")
	DBConn, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	return err
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
	dsn := "host=0.0.0.0 user=postgres password=postgres dbname=postgres port=5433"
	DBConn, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(err)
	}
}
