package models

import (
	"fmt"
	"time"
)

// Database model for item
type Item struct {
	ID    uint   `gorm:"primary_key" json:"id"`
	Title string `gorm:"not null" json:"title"`
	//TODO Placement?
	Description string    `json:"description"`
	SectionId   uint      `gorm:"not null" json:"section_id"`
	WorkspaceId uint      `gorm:"not null" json:"workspace_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
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
