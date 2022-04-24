package boardapi

import (
	"github.com/ferealqq/golang-trello-copy/server/boardapi/models"
	"github.com/palantir/stacktrace"
)

// Compile-time proof of interface implementation
var _ models.BoardStorage = (*BoardService)(nil)

type BoardService struct {
	BoardList  map[int]models.Board
	MaxBoardID int
}

// NewBoardService creates a new Carer Service with the system's database connection
func NewBoardService(list map[int]models.Board, count int) models.BoardStorage {
	return &BoardService{
		BoardList:  list,
		MaxBoardID: count,
	}
}

// ListBoards returns a list of JSON documents
func (service *BoardService) ListBoards() ([]models.Board, error) {
	var list []models.Board
	for _, v := range service.BoardList {
		list = append(list, v)
	}
	return list, nil
}

// GetBoard returns a single JSON document
func (service *BoardService) GetBoard(i int) (models.Board, error) {
	Board, ok := service.BoardList[i]
	if !ok {
		return models.Board{}, stacktrace.NewError("Failure trying to retrieve Board")
	}
	return Board, nil
}

// AddBoard adds a Board JSON document, returns the JSON document with the generated id
func (service *BoardService) AddBoard(b models.Board) (models.Board, error) {
	service.MaxBoardID = service.MaxBoardID + 1
	b.ID = service.MaxBoardID
	service.BoardList[service.MaxBoardID] = b
	return b, nil
}

// UpdateBoard updates an existing Board
func (service *BoardService) UpdateBoard(b models.Board) (models.Board, error) {
	id := b.ID
	_, ok := service.BoardList[id]
	if !ok {
		return b, stacktrace.NewError("Failure trying to update Board")
	}
	service.BoardList[id] = b
	return service.BoardList[id], nil
}

// DeleteBoard deletes a Board
func (service *BoardService) DeleteBoard(i int) error {
	_, ok := service.BoardList[i]
	if !ok {
		return stacktrace.NewError("Failure trying to delete Board")
	}
	delete(service.BoardList, i)
	return nil
}

// CreateMockDataSet initialises a database for test purposes. It returns a list of Board objects
// as well as the new max object count
func CreateMockBoardSet() (map[int]models.Board, int) {
	list := make(map[int]models.Board)
	list[0] = models.Board{
		ID:              0,
		Title:       	 "App development",
		Description:     "This is a board for app development",
	}
	list[1] = models.Board{
		ID:              1,
		Title:           "Web development",
		Description:     "This is a board for web development",
	}
	return list, len(list) - 1
}
