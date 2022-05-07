package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"testing"

	. "github.com/ferealqq/golang-trello-copy/server/boardapi/models"
	app "github.com/ferealqq/golang-trello-copy/server/pkg/appenv"
	ctrl "github.com/ferealqq/golang-trello-copy/server/pkg/controller"
	"github.com/ferealqq/golang-trello-copy/server/pkg/database"
	. "github.com/ferealqq/golang-trello-copy/server/pkg/testing"
	. "github.com/ferealqq/golang-trello-copy/server/seeders"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestListItemsHandler(t *testing.T) {
	action := HttpTestAction[Board]{
		Method: http.MethodGet,
		RouterFunc: func(e *gin.Engine, ae app.AppEnv) {
			e.GET("/items", ctrl.MakeHandler(ae, ListItemsHandler))
		},
		ReqPath: "/items",
		Seeders: []func(db *gorm.DB){SeedItems},
		Tables:  []string{"items"},
	}
	response := action.Run()
	assert.Equal(t, http.StatusOK, response.Code, "they should be equal")
	var d map[string]interface{}

	if err := json.Unmarshal(response.Body.Bytes(), &d); err != nil {
		assert.Fail(t, "Unmarshal should not fail")
		return
	}
	var count int64
	err := database.DBConn.Model(&Item{}).Count(&count)
	if err.Error != nil {
		assert.Fail(t, "Fetching count should not fail")
	}
	assert.Equal(t, int(count), int(d["count"].(float64)), "they should be equal")
}

func TestListItemsHandlerFromWorkspaceAndSection(t *testing.T) {
	action := HttpTestAction[Board]{
		Method: http.MethodGet,
		RouterFunc: func(e *gin.Engine, ae app.AppEnv) {
			e.GET("/items", ctrl.MakeHandler(ae, ListItemsHandler))
		},
		Seeders: []func(db *gorm.DB){},
		Tables:  []string{},
	}
	w := CreateWorkspaceFaker(database.DBConn).Model
	b := Board{
		Title:       "Only this boards sections will be listed",
		Description: "This is a test board",
		WorkspaceId: w.ID,
	}
	database.DBConn.Create(&b)
	s := CreateSection(database.DBConn, "Section one of results wanted", "this is a test section", b.ID).Model
	s1 := CreateSection(database.DBConn, "Section two of results wanted", "this is a test section", b.ID).Model
	for i := 0; i < 6; i++ {
		if i%2 == 0 {
			CreateItem(database.DBConn, "This should not be returned with GET /items", "item for section description", w.ID, s1.ID)
		} else {
			fmt.Println("creating item")
			CreateItem(database.DBConn, "This should be found with GET /items", "item for section description", w.ID, s.ID)
		}
	}

	// Query with multiple board ids
	action.ReqPath = "/items?WorkspaceId=" +
		strconv.FormatUint(uint64(w.ID), 10) + "&SectionId=" +
		strconv.FormatUint(uint64(s.ID), 10)

	response := action.Run()
	assert.Equal(t, http.StatusOK, response.Code, "they should be equal")
	var d map[string]interface{}

	if err := json.Unmarshal(response.Body.Bytes(), &d); err != nil {
		assert.Fail(t, "Unmarshal should not fail")
		return
	}
	assert.Equal(t, 3, int(d["count"].(float64)), "they should be equal")
}

func TestPatchItemHandler(t *testing.T) {
	// BoardId 1 is created by SeedSections function, by default if there is no board in the database SeedSections will create a new board
	i := Item{
		SectionId: 7,
	}
	b, _ := json.Marshal(i)

	action := HttpTestAction[Section]{
		Method: http.MethodPatch,
		RouterFunc: func(e *gin.Engine, ae app.AppEnv) {
			e.PATCH("/items/:id", ctrl.MakeHandler(ae, UpdateItemHandler))
		},
		ReqPath: "/items/1",
		Body:    bytes.NewReader(b),
		Seeders: []func(db *gorm.DB){SeedWorkspaces, SeedSections, SeedItems, func(db *gorm.DB) {
			// generate lots of random sections
			for i := 0; i < 5; i++ {
				CreateSectionFaker(db)
			}
		}},
		Tables: []string{"workspaces", "sections", "items"},
	}

	response := action.Run()

	assert.Equal(t, http.StatusOK, response.Code, "they should be equal")

	var rItem map[string]interface{}
	if err := json.Unmarshal(response.Body.Bytes(), &rItem); err != nil {
		assert.Fail(t, "Unmarshal should not fail")
		return
	}
	var item Item
	assert.Nil(t, database.DBConn.First(&item, 1).Error, "item should be found")
	assert.Equal(t, item.SectionId, uint(7), "they should be equal")
	assert.NotNil(t, item.Title, "should not be nil")
	assert.NotNil(t, item.Description, "should not be nil")
	assert.NotNil(t, item.WorkspaceId, "should not be nil")
	assert.NotNil(t, item.UpdatedAt, "should not be nil")
}

func TestCreateItemHandler(t *testing.T) {
	i := Item{
		Title:       "This is a new item",
		Description: "This is a description for the new item that is about to be created",
		WorkspaceId: 1,
		SectionId:   1,
	}
	b, _ := json.Marshal(i)

	action := HttpTestAction[Item]{
		Method: http.MethodPost,
		RouterFunc: func(e *gin.Engine, ae app.AppEnv) {
			e.POST("/items", ctrl.MakeHandler(ae, CreateItemHandler))
		},
		ReqPath: "/items",
		Body:    bytes.NewReader(b),
		Seeders: []func(db *gorm.DB){SeedItems},
		Tables:  []string{"items"},
	}

	response := action.Run()

	assert.Equal(t, http.StatusCreated, response.Code, "they should be equal")

	var rItem map[string]interface{}
	if err := json.Unmarshal(response.Body.Bytes(), &rItem); err != nil {
		assert.Fail(t, "Unmarshal should not fail")
		return
	}
	var item Item
	err := database.DBConn.First(&item, rItem["ID"]).Error
	if err != nil {
		assert.Fail(t, "Fetching section should not fail")
	}
}

func TestGetItemHandler(t *testing.T) {
	action := HttpTestAction[Section]{
		Method: http.MethodGet,
		RouterFunc: func(e *gin.Engine, ae app.AppEnv) {
			e.GET("/items/:id", ctrl.MakeHandler(ae, GetItemHandler))
		},
		ReqPath: "/items/1",
		Seeders: []func(db *gorm.DB){SeedItems},
		Tables:  []string{"items"},
	}
	response := action.Run()
	assert.Equal(t, http.StatusOK, response.Code, "they should be equal")
	var d map[string]interface{}

	if err := json.Unmarshal(response.Body.Bytes(), &d); err != nil {
		assert.Fail(t, "Unmarshal should not fail")
		return
	}
	var item Item
	err := database.DBConn.First(&item, d["ID"]).Error
	if err != nil {
		assert.Fail(t, "Fetching item should not fail")
	}

	assert.Equal(t, item.Title, d["Title"], "they should be equal")
	assert.Equal(t, item.Description, d["Description"], "they should be equal")
	assert.Equal(t, int(item.SectionId), int(d["SectionId"].(float64)), "they should be equal")
	assert.Equal(t, int(item.WorkspaceId), int(d["WorkspaceId"].(float64)), "they should be equal")
}

func TestDeleteItemHandler(t *testing.T) {
	action := HttpTestAction[Section]{
		Method: http.MethodDelete,
		RouterFunc: func(e *gin.Engine, ae app.AppEnv) {
			e.DELETE("/items/:id", ctrl.MakeHandler(ae, DeleteItemHandler))
		},
		ReqPath: "/items/1",
		Seeders: []func(db *gorm.DB){SeedItems},
		Tables:  []string{"items"},
	}

	response := action.Run()

	assert.Equal(t, http.StatusOK, response.Code, "they should be equal")

	var item Item
	assert.True(t, errors.Is(database.DBConn.First(&item, 1).Error, gorm.ErrRecordNotFound), "they should be equal")
}
