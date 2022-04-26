package seeders

import (
	"errors"
	"log"

	. "github.com/ferealqq/golang-trello-copy/server/boardapi/models"
	. "github.com/ferealqq/golang-trello-copy/server/pkg/seed"
	"gorm.io/gorm"
)

func CreateSection(db *gorm.DB, title string, desc string, boardId uint) (*gorm.DB) {
	return db.Create(&Section{
		Title: title,
		Description: desc,
		BoardId: boardId,
	})
}

func SectionAll() []Seed {
	return []Seed{
		{
			Name: "Three sections for board 1",
			Run: func(db *gorm.DB) (*gorm.DB) {
				var board Board
				result := db.First(&board)
				if result.Error == nil {
					for i := 0; i < 3; i++ {
						//FIXME Could be done with a batch insert but i'm lazy
						res := CreateSection(db, "Section 1", "This is a section for the board", board.ID)
						if(res.Error != nil){
							panic(res.Error)
						}
					}
					return db
				} else if errors.Is(result.Error, gorm.ErrRecordNotFound) {
					CreateBoard(db, "REST API development", "This is a board for the REST API development")
					result := db.First(&board)
					if result.Error == nil {
						for i := 0; i < 3; i++ {
							//FIXME Could be done with a batch insert but i'm lazy
							res := CreateSection(db, "Section 1", "This is a section for the board", board.ID)
							if(res.Error != nil){
								panic(res.Error)
							}
						}
						return db
					}else{
						panic(result.Error)
					}
				} else {
					panic(result.Error)
				}
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