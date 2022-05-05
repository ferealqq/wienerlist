package utils

import (
	"sync"

	"gorm.io/gorm/schema"
)

func GetTableName(Model interface{}) string {
	s, _ := schema.Parse(Model, &sync.Map{}, schema.NamingStrategy{})
	return s.Table
}

func GetColumnName(col string) string {
	return schema.NamingStrategy{}.ColumnName("", col)
}
