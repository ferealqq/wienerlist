package models

import (
	"testing"

	"github.com/ferealqq/golang-trello-copy/server/pkg/database"
	"github.com/stretchr/testify/assert"
)


func TestCreateBoard(t *testing.T) {
	database.TestDBInit()

	b := Board{
		Title:       "This is a new test board",
		Description: "This is a board for app development",
	}
	result := database.DBConn.Create(&b)
	if assert.Nil(t, result.Error) {
		assert.NotNil(t,b.ID, "should exist")
		assert.Equal(t, "This is a new test board", b.Title, "they should be equal")
		assert.Equal(t, "This is a board for app development", b.Description, "they should be equal")
	}

	defer database.Close()
}

func TestGetBoard(t *testing.T){
	database.TestDBInit()

	b := Board{
		Title:       "This is a new test board",
		Description: "This is a board for app development",
	}
	result := database.DBConn.Create(&b)
	if assert.Nil(t, result.Error) {
		assert.NotNil(t,b.ID, "should exist")
		assert.Equal(t, "This is a new test board", b.Title, "they should be equal")
		assert.Equal(t, "This is a board for app development", b.Description, "they should be equal")
	}
	bo := Board{}
	result = database.DBConn.First(&bo, 2)
	assert.Nil(t, result.Error)
	assert.NotNil(t, bo)

	defer database.Close()
}