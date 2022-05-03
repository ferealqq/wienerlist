package models

import (
	"fmt"
	"time"
)

// FIXME figure out how to use migrations in production. For example if you add something AutoMigrate will take care of it. But if you delete something from a model, it does not migrate
// Database model for section
type Section struct {
	ID          uint   `gorm:"primary_key"`
	Title       string `gorm:"not null"`
	Description string
	Placement   int
	BoardId     uint   `gorm:"not null"`
	Items       []Item `gorm:"foreignkey:SectionId"`
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
