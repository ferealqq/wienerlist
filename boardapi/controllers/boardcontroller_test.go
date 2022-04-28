package controllers

import (
	"testing"

	. "github.com/ferealqq/golang-trello-copy/server/pkg/container"
	"github.com/stretchr/testify/assert"
)

func TestMakeHandler(t *testing.T) {
	a := CreateContextForTestSetup()

	assert.NotNil(t, MakeHandler(a, GetBoardHandler))
}

/**
type TestAction struct {
	Name       string
	RouterPath string
	ReqPath    string
	Handler    func(http.ResponseWriter, *http.Request, AppEnv)
	Method     string
	Body       io.Reader
	Seeders    []func(db *gorm.DB)
	Tables     []string
}

// FIXME create test "suite" so you can use multiple database connections
func BoardHandlerAction(action TestAction) *httptest.ResponseRecorder {
	appEnv := CreateContextForTestSetup()
	//FIXME figure out a better way to give out the table name, table names could change so this is a little problematic approach
	ReinitTables(appEnv.DBConn, action.Tables, action.Seeders)
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
		Seeders:    []func(db *gorm.DB){SeedBoards},
		Tables:     []string{"boards"},
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
		Method:     http.MethodPost,
		RouterPath: "/boards",
		ReqPath:    "/boards",
		Name:       "CreateBoardHandler",
		Handler:    CreateBoardHandler,
		Body:       bytes.NewReader(b),
		Seeders:    []func(db *gorm.DB){SeedBoards},
		Tables:     []string{"boards"},
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
		Method:     http.MethodDelete,
		RouterPath: "/boards/{bid:[0-9]+}",
		ReqPath:    "/boards/1",
		Name:       "DeleteBoardHandler",
		Handler:    DeleteBoardHandler,
		Seeders:    []func(db *gorm.DB){SeedBoards},
		Tables:     []string{"boards"},
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
		Method:     http.MethodGet,
		RouterPath: "/boards/{bid:[0-9]+}",
		ReqPath:    "/boards/1",
		Name:       "GetBoardHandler",
		Handler:    GetBoardHandler,
		Seeders:    []func(db *gorm.DB){SeedBoards},
		Tables:     []string{"boards"},
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
		Method:     http.MethodPut,
		RouterPath: "/boards/{bid:[0-9]+}",
		ReqPath:    "/boards/1",
		Name:       "UpdateBoardHandler",
		Handler:    UpdateBoardHandler,
		Body:       bytes.NewReader(b),
		Seeders:    []func(db *gorm.DB){SeedBoards},
		Tables:     []string{"boards"},
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
		Method:     http.MethodGet,
		RouterPath: "/boards/{bid:[0-9]+}",
		ReqPath:    "/boards/1",
		Name:       "GetBoardHandler",
		Handler:    GetBoardHandler,
		Seeders:    []func(db *gorm.DB){CreateBoardFaker, SeedSections},
		Tables:     []string{"boards", "sections"},
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
*/
