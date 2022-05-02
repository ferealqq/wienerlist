package controllers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	. "github.com/ferealqq/golang-trello-copy/server/boardapi/models"
	"github.com/ferealqq/golang-trello-copy/server/migrations"
	app "github.com/ferealqq/golang-trello-copy/server/pkg/appenv"
	ctrl "github.com/ferealqq/golang-trello-copy/server/pkg/controller"
	"github.com/ferealqq/golang-trello-copy/server/pkg/database"
	. "github.com/ferealqq/golang-trello-copy/server/pkg/testing"
	. "github.com/ferealqq/golang-trello-copy/server/seeders"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/negroni"
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

type TestAction struct {
	Name string
	// Functionality of the router
	RouterFunc func(*gin.Engine, app.AppEnv)
	// Path to which the request will be sent
	ReqPath string
	Handler func(ctrl.BaseController[Board])
	Method  string
	Body    io.Reader
	Seeders []func(db *gorm.DB)
	Tables  []string
}

// FIXME create test "suite" so you can use multiple database connections
func BoardHandlerAction(action TestAction) *httptest.ResponseRecorder {
	appenv := app.CreateTestAppEnv()
	//FIXME figure out a better way to give out the table name, table names could change so this is a little problematic approach
	ReinitTables(database.DBConn, action.Tables, action.Seeders)
	r, _ := http.NewRequest(action.Method, action.ReqPath, action.Body)
	w := httptest.NewRecorder()
	router := gin.Default()
	// recreate route from routes.go
	action.RouterFunc(router, appenv)
	n := negroni.New()
	n.UseHandler(router)
	n.ServeHTTP(w, r)

	return w
}

func TestListBoardsHandler(t *testing.T) {
	response := BoardHandlerAction(TestAction{
		Method: http.MethodGet,
		RouterFunc: func(e *gin.Engine, ae app.AppEnv) {
			e.GET("/boards", ctrl.MakeHandler(ae, ListBoardsHandler))
		},
		ReqPath: "/boards",
		Seeders: []func(db *gorm.DB){SeedBoards},
		Tables:  []string{"boards"},
	})
	assert.Equal(t, http.StatusOK, response.Code, "they should be equal")
	var d map[string]interface{}

	if err := json.Unmarshal(response.Body.Bytes(), &d); err != nil {
		assert.Fail(t, "Unmarshal should not fail")
		return
	}
	assert.Equal(t, len(BoardAll()), int(d["count"].(float64)), "they should be equal")
}

func TestPostBoardHandler(t *testing.T) {
	// create a new board
	board := Board{
		Title:       "Test Board",
		Description: "This is a test board",
	}
	// to io reader
	b, _ := json.Marshal(board)

	response := BoardHandlerAction(TestAction{
		Method:  http.MethodPost,
		ReqPath: "/boards",
		RouterFunc: func(e *gin.Engine, ae app.AppEnv) {
			e.POST("/boards", ctrl.MakeHandler(ae, CreateBoardHandler))
		},
		Body:    bytes.NewReader(b),
		Seeders: []func(db *gorm.DB){SeedBoards},
		Tables:  []string{"boards"},
	})

	assert.Equal(t, http.StatusCreated, response.Code, "they should be equal")
	var d map[string]interface{}
	var boards []Board
	database.DBConn.Find(&boards)

	if err := json.Unmarshal(response.Body.Bytes(), &d); err != nil {
		assert.Fail(t, "Unmarshal should not fail")
		return
	}
	if assert.Equal(t, len(boards), int(d["ID"].(float64)), "they should be equal") {
		assert.Equal(t, board.Title, d["Title"], "they should be equal")
		assert.Equal(t, board.Description, d["Description"], "they should be equal")
	}
}

func TestDeleteBoardHandler(t *testing.T) {
	response := BoardHandlerAction(TestAction{
		Method: http.MethodDelete,
		RouterFunc: func(e *gin.Engine, ae app.AppEnv) {
			e.DELETE("/boards/:id", ctrl.MakeHandler(ae, DeleteBoardHandler))
		},
		ReqPath: "/boards/1",
		Seeders: []func(db *gorm.DB){SeedBoards},
		Tables:  []string{"boards"},
	})

	var boards []Board
	database.DBConn.Find(&boards)
	var d map[string]interface{}
	assert.Equal(t, http.StatusOK, response.Code, "they should be equal")

	if err := json.Unmarshal(response.Body.Bytes(), &d); err != nil {
		assert.Fail(t, "Unmarshal should not fail")
		return
	}
	assert.Equal(t, int(0), int(d["ID"].(float64)), "they should be equal")
}

func TestGetBoardHandler(t *testing.T) {
	response := BoardHandlerAction(TestAction{
		Method: http.MethodGet,
		RouterFunc: func(e *gin.Engine, ae app.AppEnv) {
			e.GET("/boards/:id", ctrl.MakeHandler(ae, GetBoardHandler))
		},
		ReqPath: "/boards/1",
		Handler: GetBoardHandler,
		Seeders: []func(db *gorm.DB){SeedBoards},
		Tables:  []string{"boards"},
	})

	var boards []Board
	database.DBConn.Find(&boards)
	var d map[string]interface{}
	assert.Equal(t, http.StatusOK, response.Code, "they should be equal")

	if err := json.Unmarshal(response.Body.Bytes(), &d); err != nil {
		assert.Fail(t, "Unmarshal should not fail")
		return
	}
	assert.Equal(t, int(1), int(d["ID"].(float64)), "they should be equal")
}

func TestUpdateBoardHandler(t *testing.T) {
	// create a new board
	board := Board{
		Title:       "Title change",
		Description: "Desc change",
	}
	// to io reader
	b, _ := json.Marshal(board)

	response := BoardHandlerAction(TestAction{
		Method:  http.MethodPut,
		ReqPath: "/boards/1",
		RouterFunc: func(e *gin.Engine, ae app.AppEnv) {
			e.PUT("/boards/:id", ctrl.MakeHandler(ae, UpdateBoardHandler))
		},
		Body:    bytes.NewReader(b),
		Seeders: []func(db *gorm.DB){SeedBoards},
		Tables:  []string{"boards"},
	})

	var boards []Board
	database.DBConn.Find(&boards)
	var d map[string]interface{}
	assert.Equal(t, http.StatusOK, response.Code, "they should be equal")

	if err := json.Unmarshal(response.Body.Bytes(), &d); err != nil {
		assert.Fail(t, "Unmarshal should not fail")
		return
	}
	assert.Equal(t, int(1), int(d["ID"].(float64)), "they should be equal")
	assert.Equal(t, board.Title, d["Title"], "they should be equal")
	assert.Equal(t, board.Description, d["Description"], "they should be equal")
}

func TestPreloadGetBoard(t *testing.T) {
	response := BoardHandlerAction(TestAction{
		Method: http.MethodGet,
		RouterFunc: func(e *gin.Engine, ae app.AppEnv) {
			e.GET("/boards/:id", ctrl.MakeHandler(ae, GetBoardHandler))
		},
		ReqPath: "/boards/1",
		Seeders: []func(db *gorm.DB){CreateBoardFaker, SeedSections},
		Tables:  []string{"boards", "sections"},
	})

	var board Board
	database.DBConn.Preload("Sections").First(&board)
	var d map[string]interface{}
	assert.Equal(t, http.StatusOK, response.Code, "they should be equal")

	if err := json.Unmarshal(response.Body.Bytes(), &d); err != nil {
		assert.Fail(t, "Unmarshal should not fail")
		return
	}
	assert.Equal(t, int(1), int(d["ID"].(float64)), "they should be equal")
	assert.Equal(t, len(board.Sections), len(d["Sections"].([]interface{})), "they should be equal")
}
