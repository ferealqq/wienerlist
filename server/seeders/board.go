package seeders

import (
	"log"

	faker "github.com/bxcodec/faker/v3"
	. "github.com/ferealqq/golang-trello-copy/server/boardapi/models"
	. "github.com/ferealqq/golang-trello-copy/server/pkg/seed"
	"gorm.io/gorm"
)

func CreateBoard(db *gorm.DB, title string, desc string) (*gorm.DB) {
	return db.Create(&Board{
		Title: title,
		Description: desc,
	})
}

func CreateBoardFaker(db *gorm.DB) {
	db.Create(&Board{
		Title: faker.Word(),
		Description: faker.Sentence(),
	});
}


func BoardAll() []Seed {
	return []Seed{
		{
			Name: "Board 1",
			Run: func(db *gorm.DB) (*gorm.DB) {
				return CreateBoard(db, "REST API development", "This is a board for the REST API development")
			},	
		},
		{
			Name: "Board 2",
			Run: func(db *gorm.DB) (*gorm.DB) {
				return CreateBoard(db, "Frontend development", "This is a board for the frontend development")
			},
		},
	}
}


func SeedBoards(db *gorm.DB) {
	for _, seed := range BoardAll() {
		if result := seed.Run(db); result.Error != nil {
			log.Fatalf("Running seed '%s', failed with error: %s", seed.Name, result.Error)
		}
	}
}