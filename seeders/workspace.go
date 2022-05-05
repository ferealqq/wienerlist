package seeders

import (
	"log"

	faker "github.com/bxcodec/faker/v3"
	"github.com/ferealqq/golang-trello-copy/server/boardapi/models"
	"github.com/ferealqq/golang-trello-copy/server/pkg/seed"
	"gorm.io/gorm"
)

func CreateWorkspace(db *gorm.DB, title string, desc string) *seed.SeedOut[models.Workspace] {
	return seed.SeedModel(db, models.Workspace{
		Title:       title,
		Description: desc,
	})
}

func CreateWorkspaceFaker(db *gorm.DB) *seed.SeedOut[models.Workspace] {
	return CreateWorkspace(db, faker.Word(), faker.Sentence())
}

func WorkspaceAll() []seed.Seed[models.Workspace] {
	return []seed.Seed[models.Workspace]{
		{
			Name: "Workspace 1",
			Run: func(db *gorm.DB) *seed.SeedOut[models.Workspace] {
				return CreateWorkspace(
					db,
					"Webso",
					"this is a workspace for the WEBSO",
				)
			},
		},
		{
			Name: "Workspace 2",
			Run: func(db *gorm.DB) *seed.SeedOut[models.Workspace] {
				return CreateWorkspace(
					db,
					"SpurdoSpärde",
					"this is a workspace for the SpurdoSpärde",
				)
			},
		},
	}
}

func SeedWorkspaces(db *gorm.DB) {
	for _, seed := range WorkspaceAll() {
		if result := seed.Run(db); result.Error != nil {
			log.Fatalf("Running seed '%s', failed with error: %s", seed.Name, result.Error)
		}
	}
}
