package seeders

import (
	"log"

	faker "github.com/bxcodec/faker/v3"
	"github.com/ferealqq/golang-trello-copy/server/boardapi/models"
	"github.com/ferealqq/golang-trello-copy/server/pkg/seed"
	"gorm.io/gorm"
)

func CreateBoard(db *gorm.DB, title string, desc string, workspace uint) *seed.SeedOut[models.Board] {
	return seed.SeedModel(db, models.Board{
		Title:       title,
		Description: desc,
		WorkspaceId: workspace,
	})
}

func CreateBoardFaker(db *gorm.DB) *seed.SeedOut[models.Board] {
	return seed.SeedModel(db, models.Board{
		Title:       faker.Word(),
		Description: faker.Sentence(),
		WorkspaceId: CreateWorkspaceFaker(db).Model.ID,
	})
}

func BoardAll() []seed.Seed[models.Board] {
	return []seed.Seed[models.Board]{
		{
			Name: "Board 1",
			Run: func(db *gorm.DB) *seed.SeedOut[models.Board] {
				return CreateBoard(
					db,
					"REST API development",
					"This is a board for the REST API development",
					1,
				)
			},
		},
		{
			Name: "Board 2",
			Run: func(db *gorm.DB) *seed.SeedOut[models.Board] {
				return CreateBoard(
					db,
					"Frontend development",
					"This is a board for the frontend development",
					1,
				)
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
