package boardapi

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/ferealqq/golang-trello-copy/server/boardapi/models"
	"github.com/ferealqq/golang-trello-copy/server/pkg/database"
	. "github.com/ferealqq/golang-trello-copy/server/pkg/testing"
	"github.com/ferealqq/golang-trello-copy/server/seeders"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/negroni"
)


type TestAction struct {
	Name string
	RouterPath string
	ReqPath string
	Handler func(http.ResponseWriter, *http.Request, AppEnv)
	Method string
	Body io.Reader
}

// FIXME create test "suite" so you can use multiple database connections

func BoardHandlerAction(action TestAction) *httptest.ResponseRecorder {
	appEnv := CreateContextForTestSetup()
	//FIXME figure out a better way to give out the table name, table names could change so this is a little problematic approach
	ReinitTable(appEnv.DBConn, "boards", seeders.SeedBoards)
	r, _ := http.NewRequest(action.Method, action.ReqPath, action.Body)
	w := httptest.NewRecorder()
	router := mux.NewRouter().StrictSlash(true)
	// recreate route from routes.go
	router.
		Methods(action.Method).
		Path(action.RouterPath).
		Name(action.Name).
		Handler(MakeHandler(appEnv, action.Handler))
	n := negroni.New()
	n.UseHandler(router)
	n.ServeHTTP(w, r)

	return w
}

func TestListBoardsHandler(t *testing.T) {
	response := BoardHandlerAction(TestAction{
		Method:     http.MethodGet,
		RouterPath: "/boards",
		ReqPath:    "/boards",
		Name:       "ListBoardsHandler",
		Handler:    ListBoardsHandler,
	})
	assert.Equal(t, http.StatusOK, response.Code, "they should be equal")
	var d map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &d)
	assert.Equal(t, len(seeders.BoardAll()), int(d["count"].(float64)), "they should be equal")
}

func TestPostBoardHandler(t *testing.T){
	// create a new board
	board := Board{
		Title: "Test Board",
		Description: "This is a test board",
	}
	// to io reader
	b, _ := json.Marshal(board)


	response := BoardHandlerAction(TestAction{
		Method:     http.MethodPost,
		RouterPath: "/boards",
		ReqPath:    "/boards",
		Name:    "CreateBoardHandler",
		Handler: CreateBoardHandler,
		Body:    bytes.NewReader(b),
	})			

	assert.Equal(t, http.StatusCreated, response.Code, "they should be equal")
	var d map[string]interface{}
	var boards []Board
	database.DBConn.Find(&boards)
	
	json.Unmarshal(response.Body.Bytes(), &d)
	if assert.Equal(t, len(boards), int(d["ID"].(float64)), "they should be equal") {
		assert.Equal(t, board.Title, d["Title"], "they should be equal")
		assert.Equal(t, board.Description, d["Description"], "they should be equal")
	}
}


func TestDeleteBoardHandler(t *testing.T){
	response := BoardHandlerAction(TestAction{
		Method:     http.MethodDelete,
		RouterPath: "/boards/{bid:[0-9]+}",
		ReqPath:    "/boards/1",
		Name:       "DeleteBoardHandler",
		Handler:    DeleteBoardHandler,
	})				

	var boards []Board
	database.DBConn.Find(&boards)
	var d map[string]interface{}
	assert.Equal(t, http.StatusOK, response.Code, "they should be equal")

	json.Unmarshal(response.Body.Bytes(), &d)
	assert.Equal(t, int(0), int(d["ID"].(float64)), "they should be equal") 
}

func TestGetBoardHandler(t *testing.T){
	response := BoardHandlerAction(TestAction{
		Method:     http.MethodGet,
		RouterPath: "/boards/{bid:[0-9]+}",
		ReqPath:    "/boards/1",
		Name:       "GetBoardHandler",
		Handler:    GetBoardHandler,
	})				
	
	var boards []Board
	database.DBConn.Find(&boards)
	var d map[string]interface{}
	assert.Equal(t, http.StatusOK, response.Code, "they should be equal")

	json.Unmarshal(response.Body.Bytes(), &d)
	assert.Equal(t, int(1), int(d["ID"].(float64)), "they should be equal") 
}

func TestUpdateBoardHandler(t *testing.T){
	// create a new board
	board := Board{
		Title: "Title change",
		Description: "Desc change",
	}
	// to io reader
	b, _ := json.Marshal(board)
	
	response := BoardHandlerAction(TestAction{
		Method:     http.MethodPut,
		RouterPath: "/boards/{bid:[0-9]+}",
		ReqPath:    "/boards/1",
		Name:       "UpdateBoardHandler",
		Handler:    UpdateBoardHandler,
		Body:	    bytes.NewReader(b),	
	})				
	
	var boards []Board
	database.DBConn.Find(&boards)
	var d map[string]interface{}
	assert.Equal(t, http.StatusOK, response.Code, "they should be equal")

	json.Unmarshal(response.Body.Bytes(), &d)
	assert.Equal(t, int(1), int(d["ID"].(float64)), "they should be equal") 
	assert.Equal(t, board.Title, d["Title"], "they should be equal") 
	assert.Equal(t, board.Description, d["Description"], "they should be equal") 
}