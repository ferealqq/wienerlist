package models

import (
	"fmt"
	"time"
)

// Database model for board
type Workspace struct {
	ID          uint   `gorm:"primary_key"`
	Title       string `gorm:"not null"`
	Description string
	Boards      []Board `gorm:"foreignkey:WorkspaceId"`
	Items       []Item  `gorm:"foreignkey:WorkspaceId"` // Needed for later usage of TODO cli feature
	CreatedAt   time.Time
	UpdatedAt   time.Time
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
