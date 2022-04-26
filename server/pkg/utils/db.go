package utils

import (
	"sync"

	"gorm.io/gorm/schema"
)

func GetTableName(Model interface{}) string {
	s, _ := schema.Parse(Model, &sync.Map{}, schema.NamingStrategy{})
	return s.Table
}
