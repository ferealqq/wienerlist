package models

import (
	"fmt"
	"time"
)

// Database model for board
type Section struct {
	ID          uint `gorm:"primary_key"`
	Title       string
	Description string
	Placement   int
	BoardId     uint
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// GoString implements the GoStringer interface so we can display the full struct during debugging
// usage: fmt.Printf("%#v", i)
// ensure that i is a pointer, so might need to do &i in some cases
func (s *Section) GoString() string {
	return fmt.Sprintf(`
{
	ID: %d,
	Title: %s,
	Description: %s,
	Placement: %d
}`,
		s.ID,
		s.Title,
		s.Description,
		s.Placement,
	)
}
