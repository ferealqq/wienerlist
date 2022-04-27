package models

import (
	"fmt"
	"time"
)

// Database model for board
type Board struct {
	ID          uint `gorm:"primary_key"`
	Title       string
	Description string
	Sections    []Section `gorm:"foreignkey:BoardId"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// GoString implements the GoStringer interface so we can display the full struct during debugging
// usage: fmt.Printf("%#v", i)
// ensure that i is a pointer, so might need to do &i in some cases
func (b *Board) GoString() string {
	return fmt.Sprintf(`
{
	ID: %d,
	Title: %s,
	Description: %s,
}`,
		b.ID,
		b.Title,
		b.Description,
	)
}
