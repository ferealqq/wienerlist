package models

import (
	"fmt"
)

type Board struct { 
	ID int `json:"id"`
	Title string `json:"title"`
	Description string `json:"Description"`
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

type BoardStorage interface {
	ListBoards() ([]Board, error)
	AddBoard(b Board) (Board, error)
	GetBoard(i int) (Board, error)
	UpdateBoard(b Board) (Board, error)
	DeleteBoard(i int) error
}
