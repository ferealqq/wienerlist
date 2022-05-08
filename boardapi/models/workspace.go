package models

import (
	"fmt"
	"time"
)

// Database model for board
type Workspace struct {
	ID          uint      `gorm:"primary_key" json:"id"`
	Title       string    `gorm:"not null" json:"title"`
	Description string    `json:"description"`
	Boards      []Board   `gorm:"foreignkey:WorkspaceId" json:"boards"`
	Items       []Item    `gorm:"foreignkey:WorkspaceId" json:"items"` // Needed for later usage of TODO cli feature
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// GoString implements the GoStringer interface so we can display the full struct during debugging
// usage: fmt.Printf("%#v", i)
// ensure that i is a pointer, so might need to do &i in some cases
func (w *Workspace) GoString() string {
	return fmt.Sprintf(`
{
	ID: %d,
	Title: %s,
	Description: %s,
}`,
		w.ID,
		w.Title,
		w.Description,
	)
}
