package boardapi

import (
	"testing"

	"github.com/ferealqq/golang-trello-copy/server/boardapi/models"
	"github.com/stretchr/testify/assert"
)

func TestListBoards(t *testing.T) {
	appEnv := CreateContextForTestSetup()
	list, _ := appEnv.BoardStore.ListBoards()
	count := len(list)
	assert.Equal(t, 2, count, "There should be 2 items in the list.")
}

func TestGetBoard(t *testing.T) {
	appEnv := CreateContextForTestSetup()
	b, err := appEnv.BoardStore.GetBoard(0)
	if assert.Nil(t, err) {
		assert.Equal(t, 0, b.ID, "they should be equal")
		assert.Equal(t, "App development", b.Title, "they should be equal")
		assert.Equal(t, "This is a board for app development", b.Description, "they should be equal")
	}
}

func TestGetBoardFail(t *testing.T) {
	appEnv := CreateContextForTestSetup()
	_, err := appEnv.BoardStore.GetBoard(10)
	assert.NotNil(t, err)
}


func TestAddBoard(t *testing.T) {
	appEnv := CreateContextForTestSetup()
	b := models.Board{
		Title:       "This is a new test board",
		Description: "This is a board for app development",
	}
	b, _ = appEnv.BoardStore.AddBoard(b)
	// we should now have a board object with a database Id
	assert.Equal(t, 2, b.ID, "Expected database Id should be 2.")
	// we should now have 3 items in the list
	list, _ := appEnv.BoardStore.ListBoards()
	count := len(list)
	assert.Equal(t, 3, count, "There should be 3 items in the list.")
}

func TestUpdateBoardSuccess(t *testing.T) {
	appEnv := CreateContextForTestSetup()
	b := models.Board{
		ID: 0,
		Title: "This is a new test board",
		Description: "This is a board for app development",
	}
	b, _ = appEnv.BoardStore.UpdateBoard(b)
	// we should now have a board object with a database Id
	assert.Equal(t, 0, b.ID, "they should be equal")
	assert.Equal(t, "This is a new test board", b.Title, "they should be equal")
	assert.Equal(t, "This is a board for app development", b.Description, "they should be equal")
}

func TestUpdateBoardFail(t *testing.T) {
	appEnv := CreateContextForTestSetup()
	b := models.Board{
		ID: 100,
		Title: "This is a new test board",
		Description: "This is a board for app development",
	}
	_, err := appEnv.BoardStore.UpdateBoard(b)
	assert.NotNil(t, err)
}

func TestDeleteBoardSuccess(t *testing.T) {
	appEnv := CreateContextForTestSetup()
	err := appEnv.BoardStore.DeleteBoard(1)
	assert.Nil(t, err)
}

func TestDeleteBoardFail(t *testing.T) {
	appEnv := CreateContextForTestSetup()
	err := appEnv.BoardStore.DeleteBoard(10)
	assert.NotNil(t, err)
}