package models

import (
	"fmt"
	"time"
)

// Database model for item
type Item struct {
	ID          uint   `gorm:"primary_key"`
	Title       string `gorm:"not null"`
	Description string
	SectionId   uint `gorm:"not null"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// GoString implements the GoStringer interface so we can display the full struct during debugging
// usage: fmt.Printf("%#v", i)
// ensure that i is a pointer, so might need to do &i in some cases
func (i *Item) GoString() string {
	return fmt.Sprintf(`
{
	ID: %d,
	Title: %s,
	Description: %s,
}`,
		i.ID,
		i.Title,
		i.Description,
	)
}
