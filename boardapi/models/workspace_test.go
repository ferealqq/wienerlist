package models

import (
	"os"
	"testing"

	"github.com/ferealqq/wienerlist/pkg/database"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	database.TestDBInit()
	if err := database.DBConn.AutoMigrate(&Workspace{}, &Workspace{}, &Section{}, &Item{}); err != nil {
		panic(err)
	}

	m.Run()
	defer database.Close()
	os.Exit(0)
}

func TestCreateWorkspace(t *testing.T) {
	// FIXME randomly tests fail
	database.TestDBInit()

	b := Workspace{
		Title:       "This is a new test Workspace",
		Description: "This is a Workspace for app development",
	}
	result := database.DBConn.Create(&b)
	if assert.Nil(t, result.Error) {
		assert.NotNil(t, b.ID, "should exist")
		assert.Equal(t, "This is a new test Workspace", b.Title, "they should be equal")
		assert.Equal(t, "This is a Workspace for app development", b.Description, "they should be equal")
	}

	defer database.Close()
}

func TestGetWorkspace(t *testing.T) {
	database.TestDBInit()

	b := Workspace{
		Title:       "This is a new test Workspace",
		Description: "This is a Workspace for app development",
	}
	result := database.DBConn.Create(&b)
	if assert.Nil(t, result.Error) {
		assert.NotNil(t, b.ID, "should exist")
		assert.Equal(t, "This is a new test Workspace", b.Title, "they should be equal")
		assert.Equal(t, "This is a Workspace for app development", b.Description, "they should be equal")
	}
	bo := Workspace{}
	result = database.DBConn.First(&bo, 2)
	assert.Nil(t, result.Error)
	assert.NotNil(t, bo)

	defer database.Close()
}
