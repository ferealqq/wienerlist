package seed

import (
	"gorm.io/gorm"
)

type Seed[M interface{}] struct {
	Name string
	Run  func(*gorm.DB) *SeedOut[M]
}

type SeedOut[M interface{}] struct {
	Error error
	Model M
}

func SeedModel[M interface{}](db *gorm.DB, model M) *SeedOut[M] {
	return &SeedOut[M]{
		Error: db.Create(&model).Error,
		Model: model,
	}
}
