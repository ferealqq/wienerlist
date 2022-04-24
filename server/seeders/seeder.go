package seeders

import (
	"log"

	"github.com/ferealqq/golang-trello-copy/server/pkg/database"
)

func main() {
	database.InitDB()

	for _, seed := range BoardAll() {
		if result := seed.Run(database.DBConn); result.Error != nil {
			log.Fatalf("Running seed '%s', failed with error: %s", seed.Name, result.Error)
		}
	}
}