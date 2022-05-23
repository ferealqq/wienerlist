package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
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

func TestListWorkspacesHandler(t *testing.T) {
	action := HttpTestAction[Board]{
		Method: http.MethodGet,
		RouterFunc: func(e *gin.Engine, ae app.AppEnv) {
			e.GET("/workspaces", ctrl.MakeHandler(ae, ListWorkspacesHandler))
		},
		ReqPath: "/workspaces",
		Seeders: []func(db *gorm.DB){SeedWorkspaces},
		Tables:  []string{"workspaces"},
	}
	response := action.Run()
	assert.Equal(t, http.StatusOK, response.Code, "they should be equal")
	var d map[string]interface{}

	if err := json.Unmarshal(response.Body.Bytes(), &d); err != nil {
		assert.Fail(t, "Unmarshal should not fail")
		return
	}
	assert.Equal(t, len(WorkspaceAll()), int(d["count"].(float64)), "they should be equal")
}

func TestListWorkspacesHandlerLimit(t *testing.T) {
	action := HttpTestAction[Board]{
		Method: http.MethodGet,
		RouterFunc: func(e *gin.Engine, ae app.AppEnv) {
			e.GET("/workspaces", ctrl.MakeHandler(ae, ListWorkspacesHandler))
		},
		ReqPath: "/workspaces?limit=1",
		Seeders: []func(db *gorm.DB){SeedWorkspaces},
		Tables:  []string{"Workspaces"},
	}
	response := action.Run()
	assert.Equal(t, http.StatusOK, response.Code, "they should be equal")
	var d map[string]interface{}

	if err := json.Unmarshal(response.Body.Bytes(), &d); err != nil {
		assert.Fail(t, "Unmarshal should not fail")
		return
	}
	assert.Equal(t, 1, int(d["count"].(float64)), "they should be equal")
}

func TestPostWorkspaceHandler(t *testing.T) {
	// create a new Workspace
	space := Workspace{
		Title:       "Test Workspace",
		Description: "This is a test Workspace",
	}
	b, _ := json.Marshal(space)

	action := HttpTestAction[Workspace]{
		Method:  http.MethodPost,
		ReqPath: "/workspaces",
		RouterFunc: func(e *gin.Engine, ae app.AppEnv) {
			e.POST("/workspaces", ctrl.MakeHandler(ae, CreateWorkspaceHandler))
		},
		Body:    bytes.NewReader(b),
		Seeders: []func(db *gorm.DB){},
		Tables:  []string{"workspaces"},
	}

	response := action.Run()

	assert.Equal(t, http.StatusCreated, response.Code, "they should be equal")
	var d map[string]interface{}
	var wspace Workspace
	database.DBConn.First(&wspace)

	if err := json.Unmarshal(response.Body.Bytes(), &d); err != nil {
		assert.Fail(t, "Unmarshal should not fail")
		return
	}
	if assert.Equal(t, wspace.ID, uint(d["id"].(float64)), "they should be equal") {
		assert.Equal(t, wspace.Title, d["title"], "they should be equal")
		assert.Equal(t, wspace.Description, d["description"], "they should be equal")
	}
}

func TestDeleteWorkspaceHandler(t *testing.T) {
	action := HttpTestAction[Workspace]{
		Method: http.MethodDelete,
		RouterFunc: func(e *gin.Engine, ae app.AppEnv) {
			e.DELETE("/workspace/:id", ctrl.MakeHandler(ae, DeleteWorkspaceHandler))
		},
		ReqPath: "/workspace/1",
		Seeders: []func(db *gorm.DB){SeedWorkspaces},
		Tables:  []string{"workspaces"},
	}

	response := action.Run()

	assert.True(
		t,
		errors.Is(database.DBConn.First(&Workspace{}, 1).Error, gorm.ErrRecordNotFound),
		"Should not find the first Workspace object from the database",
	)
	var d map[string]interface{}
	assert.Equal(t, http.StatusOK, response.Code, "they should be equal")

	if err := json.Unmarshal(response.Body.Bytes(), &d); err != nil {
		assert.Fail(t, "Unmarshal should not fail")
		return
	}
	assert.Equal(t, int(0), int(d["id"].(float64)), "they should be equal")
}

func TestGetWorkspaceHandler(t *testing.T) {
	action := HttpTestAction[Workspace]{
		Method: http.MethodGet,
		RouterFunc: func(e *gin.Engine, ae app.AppEnv) {
			e.GET("/workspaces/:id", ctrl.MakeHandler(ae, GetWorkspaceHandler))
		},
		ReqPath: "/workspaces/1",
		Handler: GetWorkspaceHandler,
		Seeders: []func(db *gorm.DB){SeedWorkspaces},
		Tables:  []string{"workspaces"},
	}
	response := action.Run()

	var wspace Workspace
	assert.Nil(t, database.DBConn.First(&wspace, 1).Error, "Should not return error")
	var d map[string]interface{}
	assert.Equal(t, http.StatusOK, response.Code, "they should be equal")

	if err := json.Unmarshal(response.Body.Bytes(), &d); err != nil {
		assert.Fail(t, "Unmarshal should not fail")
		return
	}
	assert.Equal(t, int(1), int(d["id"].(float64)), "they should be equal")
}

func TestUpdateWorkspaceHandler(t *testing.T) {
	// create a new board
	wspace := Workspace{
		Title: "Title change",
	}
	// to io reader
	wJson, _ := json.Marshal(wspace)

	action := HttpTestAction[Workspace]{
		Method:  http.MethodPatch,
		ReqPath: "/workspaces/1",
		RouterFunc: func(e *gin.Engine, ae app.AppEnv) {
			e.PATCH("/workspaces/:id", ctrl.MakeHandler(ae, UpdateWorkspaceHandler))
		},
		Body:    bytes.NewReader(wJson),
		Seeders: []func(db *gorm.DB){SeedWorkspaces},
		Tables:  []string{"workspaces"},
	}
	response := action.Run()

	var w Workspace
	assert.Nil(t, database.DBConn.First(&w, 1).Error, "Should not return error")
	var d map[string]interface{}
	assert.Equal(t, http.StatusOK, response.Code, "they should be equal")

	if err := json.Unmarshal(response.Body.Bytes(), &d); err != nil {
		assert.Fail(t, "Unmarshal should not fail")
		return
	}
	assert.Equal(t, w.ID, uint(1), "they should be equal")
	assert.Equal(t, w.Title, wspace.Title, "they should be equal")
	assert.NotNil(t, w.Description, "they should be equal")
	assert.NotNil(t, w.CreatedAt, "they should not be nil")
}

func TestPreloadGetWorkspace(t *testing.T) {
	action := HttpTestAction[Board]{
		Method: http.MethodGet,
		RouterFunc: func(e *gin.Engine, ae app.AppEnv) {
			e.GET("/workspaces/:id", ctrl.MakeHandler(ae, GetWorkspaceHandler))
		},
		ReqPath: "/workspaces/1",
		Seeders: []func(db *gorm.DB){func(db *gorm.DB) {
			ws := CreateWorkspaceFaker(db).Model
			b := CreateBoard(db, "title", "desc", ws.ID).Model
			CreateBoard(db, "title", "desc", ws.ID)
			s := CreateSection(db, "title", "desc", b.ID).Model
			for i := 0; i < 3; i++ {
				CreateItem(db, "title 1", "desc 1", ws.ID, s.ID)
			}
		}},
		Tables: []string{"workspaces", "boards"},
	}
	response := action.Run()

	var w Workspace
	assert.Nil(t, database.DBConn.Preload("Boards").First(&w, 1).Error, "should not throw error")
	var d map[string]interface{}
	assert.Equal(t, http.StatusOK, response.Code, "they should be equal")

	if err := json.Unmarshal(response.Body.Bytes(), &d); err != nil {
		assert.Fail(t, "Unmarshal should not fail")
		return
	}
	assert.Equal(t, int(1), int(d["id"].(float64)), "they should be equal")
	assert.Equal(t, len(w.Boards), len(d["boards"].([]interface{})), "they should be equal")
}
