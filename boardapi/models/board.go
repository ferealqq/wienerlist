package models

import (
	"fmt"
	"time"
)

// Database model for board
type Board struct {
	ID          uint      `gorm:"primary_key" json:"id"`
	Title       string    `gorm:"not null" json:"title"`
	Description string    `json:"description"`
	WorkspaceId uint      `gorm:"not null" json:"workspace_id"`
	Sections    []Section `gorm:"foreignkey:BoardId" json:"sections"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
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
