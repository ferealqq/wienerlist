package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"testing"

	. "github.com/ferealqq/wienerlist/boardapi/models"
	"github.com/ferealqq/wienerlist/migrations"
	app "github.com/ferealqq/wienerlist/pkg/appenv"
	ctrl "github.com/ferealqq/wienerlist/pkg/controller"
	"github.com/ferealqq/wienerlist/pkg/database"
	. "github.com/ferealqq/wienerlist/pkg/testing"
	. "github.com/ferealqq/wienerlist/seeders"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestMain(m *testing.M) {
	database.TestDBInit()
	// stop execution if migrations fail because tests won't be able to run
	if err := migrations.Migrate(database.DBConn); err != nil {
		panic(err)
	}
	exitVal := m.Run()
	defer database.Close()
	os.Exit(exitVal)
}

func TestMakeHandler(t *testing.T) {
	appenv := app.CreateTestAppEnv()
	assert.NotNil(t, ctrl.MakeHandler(appenv, GetBoardHandler))
}

func TestListBoardsHandler(t *testing.T) {
	action := HttpTestAction[Board]{
		Method: http.MethodGet,
		RouterFunc: func(e *gin.Engine, ae app.AppEnv) {
			e.GET("/boards", ctrl.MakeHandler(ae, ListBoardsHandler))
		},
		ReqPath: "/boards",
		Seeders: []func(db *gorm.DB){SeedWorkspaces, SeedBoards},
		Tables:  []string{"workspaces", "boards"},
	}
	response := action.Run()
	assert.Equal(t, http.StatusOK, response.Code, "they should be equal")
	var d map[string]interface{}

	if err := json.Unmarshal(response.Body.Bytes(), &d); err != nil {
		assert.Fail(t, "Unmarshal should not fail")
		return
	}
	assert.Equal(t, len(BoardAll()), int(d["count"].(float64)), "they should be equal")
}

func TestListBoardsHandlerLimit(t *testing.T) {
	action := HttpTestAction[Board]{
		Method: http.MethodGet,
		RouterFunc: func(e *gin.Engine, ae app.AppEnv) {
			e.GET("/boards", ctrl.MakeHandler(ae, ListBoardsHandler))
		},
		ReqPath: "/boards?limit=1",
		Seeders: []func(db *gorm.DB){SeedBoards},
		Tables:  []string{"boards"},
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

func TestPostBoardHandler(t *testing.T) {
	// create a new board
	board := Board{
		Title:       "Test Board",
		Description: "This is a test board",
		WorkspaceId: 1,
	}
	b, _ := json.Marshal(board)

	action := HttpTestAction[Board]{
		Method:  http.MethodPost,
		ReqPath: "/boards",
		RouterFunc: func(e *gin.Engine, ae app.AppEnv) {
			e.POST("/boards", ctrl.MakeHandler(ae, CreateBoardHandler))
		},
		Body:    bytes.NewReader(b),
		Seeders: []func(db *gorm.DB){SeedWorkspaces, SeedBoards},
		Tables:  []string{"boards", "workspaces"},
	}

	response := action.Run()

	assert.Equal(t, http.StatusCreated, response.Code, "they should be equal")
	var d map[string]interface{}
	var boards []Board
	database.DBConn.Find(&boards)

	if err := json.Unmarshal(response.Body.Bytes(), &d); err != nil {
		assert.Fail(t, "Unmarshal should not fail")
		return
	}
	if assert.Equal(t, len(boards), int(d["id"].(float64)), "they should be equal") {
		assert.Equal(t, board.Title, d["title"], "they should be equal")
		assert.Equal(t, board.Description, d["description"], "they should be equal")
	}
}

func TestDeleteBoardHandler(t *testing.T) {
	action := HttpTestAction[Board]{
		Method: http.MethodDelete,
		RouterFunc: func(e *gin.Engine, ae app.AppEnv) {
			e.DELETE("/boards/:id", ctrl.MakeHandler(ae, DeleteBoardHandler))
		},
		ReqPath: "/boards/1",
		Seeders: []func(db *gorm.DB){SeedBoards},
		Tables:  []string{"boards"},
	}

	response := action.Run()

	var boards []Board
	database.DBConn.Find(&boards)
	var d map[string]interface{}
	assert.Equal(t, http.StatusOK, response.Code, "they should be equal")

	if err := json.Unmarshal(response.Body.Bytes(), &d); err != nil {
		assert.Fail(t, "Unmarshal should not fail")
		return
	}
	assert.Equal(t, int(0), int(d["id"].(float64)), "they should be equal")
}

func TestGetBoardHandler(t *testing.T) {
	action := HttpTestAction[Board]{
		Method: http.MethodGet,
		RouterFunc: func(e *gin.Engine, ae app.AppEnv) {
			e.GET("/boards/:id", ctrl.MakeHandler(ae, GetBoardHandler))
		},
		ReqPath: "/boards/1",
		Handler: GetBoardHandler,
		Seeders: []func(db *gorm.DB){SeedBoards},
		Tables:  []string{"boards"},
	}
	response := action.Run()

	var boards []Board
	database.DBConn.Find(&boards)
	var d map[string]interface{}
	assert.Equal(t, http.StatusOK, response.Code, "they should be equal")

	if err := json.Unmarshal(response.Body.Bytes(), &d); err != nil {
		assert.Fail(t, "Unmarshal should not fail")
		return
	}
	assert.Equal(t, int(1), int(d["id"].(float64)), "they should be equal")
}

func TestUpdateBoardHandler(t *testing.T) {
	// create a new board
	b := Board{
		Title:       "Title change",
		Description: "Desc change",
		WorkspaceId: 1,
	}
	// to io reader
	bJson, _ := json.Marshal(b)

	action := HttpTestAction[Board]{
		Method:  http.MethodPatch,
		ReqPath: "/boards/1",
		RouterFunc: func(e *gin.Engine, ae app.AppEnv) {
			e.PATCH("/boards/:id", ctrl.MakeHandler(ae, UpdateBoardHandler))
		},
		Body: bytes.NewReader(bJson),
		Seeders: []func(db *gorm.DB){func(db *gorm.DB) {
			CreateWorkspaceFaker(db)
		}, SeedBoards},
		Tables: []string{"boards"},
	}
	response := action.Run()

	var board Board
	database.DBConn.First(&board, 1)
	var d map[string]interface{}
	assert.Equal(t, http.StatusOK, response.Code, "they should be equal")

	if err := json.Unmarshal(response.Body.Bytes(), &d); err != nil {
		assert.Fail(t, "Unmarshal should not fail")
		return
	}
	assert.Equal(t, board.ID, uint(1), "they should be equal")
	assert.Equal(t, board.Title, b.Title, "they should be equal")
	assert.Equal(t, board.Description, b.Description, "they should be equal")
	assert.NotNil(t, board.CreatedAt, "they should not be nil")
	assert.NotNil(t, board.WorkspaceId, "they should not be nil")

}

func TestPreloadGetBoard(t *testing.T) {
	action := HttpTestAction[Board]{
		Method: http.MethodGet,
		RouterFunc: func(e *gin.Engine, ae app.AppEnv) {
			e.GET("/boards/:id", ctrl.MakeHandler(ae, GetBoardHandler))
		},
		ReqPath: "/boards/1",
		Seeders: []func(db *gorm.DB){func(db *gorm.DB) {
			CreateBoardFaker(db)
		}, SeedSections},
		Tables: []string{"boards", "sections"},
	}
	response := action.Run()

	var board Board
	database.DBConn.Preload("Sections").First(&board)
	var d map[string]interface{}
	assert.Equal(t, http.StatusOK, response.Code, "they should be equal")

	if err := json.Unmarshal(response.Body.Bytes(), &d); err != nil {
		assert.Fail(t, "Unmarshal should not fail")
		return
	}
	assert.Equal(t, int(1), int(d["id"].(float64)), "they should be equal")
	assert.Equal(t, len(board.Sections), len(d["sections"].([]interface{})), "they should be equal")
}
