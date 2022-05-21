package seeders

import (
	"log"

	faker "github.com/bxcodec/faker/v3"
	"github.com/ferealqq/golang-trello-copy/server/boardapi/models"
	"github.com/ferealqq/golang-trello-copy/server/pkg/seed"
	"gorm.io/gorm"
)

func CreateSection(db *gorm.DB, title string, desc string, BoardId uint) *seed.SeedOut[models.Section] {
	return seed.SeedModel(db, models.Section{
		Title:       title,
		Description: desc,
		BoardId:     BoardId,
	})
}

func CreateSectionFaker(db *gorm.DB) *seed.SeedOut[models.Section] {
	return CreateSection(db, faker.Word(), faker.Sentence(), CreateBoardFaker(db).Model.ID)
}

func CreateDefaultSections(db *gorm.DB, board uint) []models.Section {
	var sections []models.Section
	// FIXME there is no error handling at this point
	sections = append(sections, CreateSection(db, "TODO", "task to do", board).Model)
	sections = append(sections, CreateSection(db, "In Progress", "tasks that are in progress", board).Model)
	sections = append(sections, CreateSection(db, "DONE", "tasks that have been completed", board).Model)
	return sections
}

func SectionAll() []seed.Seed[[]models.Section] {
	return []seed.Seed[[]models.Section]{
		{
			Name: "Three sections for board 1",
			Run: func(db *gorm.DB) *seed.SeedOut[[]models.Section] {
				var sections []models.Section
				var board models.Board
				if err := db.First(&board).Error; err != nil {
					bResult := CreateBoardFaker(db)
					if bResult.Error != nil {
						panic(bResult.Error)
					} else {
						sections = CreateDefaultSections(db, bResult.Model.ID)
					}
				} else {
					sections = CreateDefaultSections(db, board.ID)
				}
				return &seed.SeedOut[[]models.Section]{Error: nil, Model: sections}
			},
		},
	}
}

func SeedSections(db *gorm.DB) {
	for _, seed := range SectionAll() {
		if result := seed.Run(db); result.Error != nil {
			log.Fatalf("Running seed '%s', failed with error: %s", seed.Name, result.Error)
		}
	}
}
