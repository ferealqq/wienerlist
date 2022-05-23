package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"testing"

	. "github.com/ferealqq/wienerlist/boardapi/models"
	app "github.com/ferealqq/wienerlist/pkg/appenv"
	ctrl "github.com/ferealqq/wienerlist/pkg/controller"
	"github.com/ferealqq/wienerlist/pkg/database"
	. "github.com/ferealqq/wienerlist/pkg/testing"
	. "github.com/ferealqq/wienerlist/seeders"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestListSectionsHandler(t *testing.T) {
	action := HttpTestAction[Board]{
		Method: http.MethodGet,
		RouterFunc: func(e *gin.Engine, ae app.AppEnv) {
			e.GET("/sections", ctrl.MakeHandler(ae, ListSectionsHandler))
		},
		ReqPath: "/sections",
		Seeders: []func(db *gorm.DB){SeedSections},
		Tables:  []string{"sections"},
	}
	response := action.Run()
	assert.Equal(t, http.StatusOK, response.Code, "they should be equal")
	var d map[string]interface{}

	if err := json.Unmarshal(response.Body.Bytes(), &d); err != nil {
		assert.Fail(t, "Unmarshal should not fail")
		return
	}
	var count int64
	err := database.DBConn.Model(&Section{}).Count(&count)
	if err.Error != nil {
		assert.Fail(t, "Fetching count should not fail")
	}
	assert.Equal(t, int(count), int(d["count"].(float64)), "they should be equal")
}

func TestListSectionsHandlerFromBoard(t *testing.T) {
	action := HttpTestAction[Board]{
		Method: http.MethodGet,
		RouterFunc: func(e *gin.Engine, ae app.AppEnv) {
			e.GET("/sections", ctrl.MakeHandler(ae, ListSectionsHandler))
		},
		Seeders: []func(db *gorm.DB){SeedSections},
		Tables:  []string{},
	}
	w := CreateWorkspaceFaker(database.DBConn).Model
	b1 := Board{
		Title:       "Only this boards sections will be listed",
		Description: "This is a test board",
		WorkspaceId: w.ID,
	}
	database.DBConn.Create(&b1)
	CreateSection(database.DBConn, "Section one of results wanted", "this is a test section", b1.ID)
	CreateSection(database.DBConn, "Section two of results wanted", "this is a test section", b1.ID)
	b2 := Board{
		Title:       "Only this boards sections will be listed",
		Description: "This is a test board",
		WorkspaceId: w.ID,
	}
	database.DBConn.Create(&b2)
	CreateSection(database.DBConn, "Section one of results wanted", "this is a test section", b2.ID)

	// Query with multiple board ids
	action.ReqPath = "/sections?board_id=" +
		strconv.FormatUint(uint64(b1.ID), 10) + "&board_id=" +
		strconv.FormatUint(uint64(b2.ID), 10)

	response := action.Run()
	assert.Equal(t, http.StatusOK, response.Code, "they should be equal")
	var d map[string]interface{}

	if err := json.Unmarshal(response.Body.Bytes(), &d); err != nil {
		assert.Fail(t, "Unmarshal should not fail")
		return
	}
	assert.Equal(t, 3, int(d["count"].(float64)), "they should be equal")
}

func TestCreateSectionHandler(t *testing.T) {
	// BoardId 1 is created by SeedSections function, by default if there is no board in the database SeedSections will create a new board
	s := Section{
		Title:       "Test section",
		Description: "This is a test section",
		BoardId:     1,
	}
	b, _ := json.Marshal(s)

	action := HttpTestAction[Section]{
		Method: http.MethodPost,
		RouterFunc: func(e *gin.Engine, ae app.AppEnv) {
			e.POST("/sections", ctrl.MakeHandler(ae, CreateSectionHandler))
		},
		ReqPath: "/sections",
		Body:    bytes.NewReader(b),
		Seeders: []func(db *gorm.DB){SeedSections},
		Tables:  []string{"sections"},
	}

	response := action.Run()

	assert.Equal(t, http.StatusCreated, response.Code, "they should be equal")

	var responseSection map[string]interface{}
	if err := json.Unmarshal(response.Body.Bytes(), &responseSection); err != nil {
		assert.Fail(t, "Unmarshal should not fail")
		return
	}
	var section Section
	err := database.DBConn.First(&section, responseSection["id"]).Error
	if err != nil {
		assert.Fail(t, "Fetching section should not fail")
	}
}

func TestGetSectionHandler(t *testing.T) {
	action := HttpTestAction[Section]{
		Method: http.MethodGet,
		RouterFunc: func(e *gin.Engine, ae app.AppEnv) {
			e.GET("/sections/:id", ctrl.MakeHandler(ae, GetSectionHandler))
		},
		ReqPath: "/sections/1",
		Seeders: []func(db *gorm.DB){SeedSections},
		Tables:  []string{"sections"},
	}
	response := action.Run()
	assert.Equal(t, http.StatusOK, response.Code, "they should be equal")
	var d map[string]interface{}

	if err := json.Unmarshal(response.Body.Bytes(), &d); err != nil {
		assert.Fail(t, "Unmarshal should not fail")
		return
	}
	var section Section
	err := database.DBConn.First(&section, d["id"]).Error
	if err != nil {
		assert.Fail(t, "Fetching section should not fail")
	}

	assert.Equal(t, section.Title, d["title"], "they should be equal")
	assert.Equal(t, section.Description, d["description"], "they should be equal")
	assert.Equal(t, int(section.BoardId), int(d["board_id"].(float64)), "they should be equal")
}

func TestUpdateSectionHandler(t *testing.T) {
	s := Section{
		Title:       "Update test section",
		Description: "Updated description",
		Placement:   3,
	}

	b, _ := json.Marshal(s)

	action := HttpTestAction[Section]{
		Method: http.MethodPatch,
		RouterFunc: func(e *gin.Engine, ae app.AppEnv) {
			e.PATCH("/sections/:id", ctrl.MakeHandler(ae, UpdateSectionHandler))
		},
		ReqPath: "/sections/1",
		Body:    bytes.NewReader(b),
		Seeders: []func(db *gorm.DB){SeedSections},
		Tables:  []string{"sections"},
	}

	response := action.Run()

	assert.Equal(t, http.StatusOK, response.Code, "they should be equal")

	var responseSection map[string]interface{}
	if err := json.Unmarshal(response.Body.Bytes(), &responseSection); err != nil {
		assert.Fail(t, "Unmarshal should not fail")
	}

	var section Section

	assert.Nil(t, database.DBConn.First(&section, 1).Error)
	assert.Equal(t, section.Title, s.Title, "they should be equal")
	assert.Equal(t, section.Description, s.Description, "they should be equal")
	assert.Equal(t, section.Placement, s.Placement, "they should be equal")
	assert.NotNil(t, section.BoardId, "they should not be nil")
}

func TestDeleteSectionHandler(t *testing.T) {
	action := HttpTestAction[Section]{
		Method: http.MethodDelete,
		RouterFunc: func(e *gin.Engine, ae app.AppEnv) {
			e.DELETE("/sections/:id", ctrl.MakeHandler(ae, DeleteSectionHandler))
		},
		ReqPath: "/sections/1",
		Seeders: []func(db *gorm.DB){SeedSections},
		Tables:  []string{"sections"},
	}

	response := action.Run()

	assert.Equal(t, http.StatusOK, response.Code, "they should be equal")

	var section Section
	assert.True(t, errors.Is(database.DBConn.First(&section, 1).Error, gorm.ErrRecordNotFound), "they should be equal")
}
