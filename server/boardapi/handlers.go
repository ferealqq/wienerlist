package boardapi

import (
	"encoding/json"
	"net/http"
	"strconv"

	models "github.com/ferealqq/golang-trello-copy/server/boardapi/models"
	"github.com/ferealqq/golang-trello-copy/server/pkg/health"
	"github.com/ferealqq/golang-trello-copy/server/pkg/status"

	// . "github.com/ferealqq/golang-trello-copy/server/pkg/utils"
	"github.com/gorilla/mux"
)

// HandlerFunc is a custom implementation of the http.HandlerFunc
type HandlerFunc func(http.ResponseWriter, *http.Request, AppEnv)

// MakeHandler allows us to pass an environment struct to our handlers, without resorting to global
// variables. It accepts an environment (Env) struct and our own handler function. It returns
// a function of the type http.HandlerFunc so can be passed on to the HandlerFunc in main.go.
func MakeHandler(appEnv AppEnv, fn func(http.ResponseWriter, *http.Request, AppEnv)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Terry Pratchett tribute
		w.Header().Set("X-Clacks-Overhead", "GNU Terry Pratchett")
		// return function with AppEnv
		fn(w, r, appEnv)
	}
}

// HealthcheckHandler returns useful info about the app
func HealthcheckHandler(w http.ResponseWriter, req *http.Request, appEnv AppEnv) {
	check := health.Check{
		AppName: "golang-trello-copy",
		Version: appEnv.Version,
	}
	appEnv.sendJSON(w, http.StatusOK, check)
}

func ListBoardsHandler(w http.ResponseWriter, req *http.Request, appEnv AppEnv) {
	var boards []models.Board
	result := appEnv.DBConn.Preload("Sections").Find(&boards)
	if result.Error != nil {
		appEnv.sendInternalServerError(w, "Error listing boards", result.Error)
		return
	}
	responseObject := make(map[string]interface{})
	responseObject["bords"] = boards
	responseObject["count"] = len(boards)
	appEnv.sendJSON(w, http.StatusOK, responseObject)
}

func CreateBoardHandler(w http.ResponseWriter, req *http.Request, appEnv AppEnv) {
	decoder := json.NewDecoder(req.Body)
	var b models.Board
	err := decoder.Decode(&b)
	if err != nil {
		appEnv.sendJSON(w, http.StatusBadRequest, status.Response{
			Status:  strconv.Itoa(http.StatusBadRequest),
			Message: "malformed board object",
		})
		return
	}
	board := models.Board{
		Title:       b.Title,
		Description: b.Description,
	}
	result := appEnv.DBConn.Create(&board)
	if result.Error != nil {
		appEnv.sendInternalServerError(w, "Error creating a board", result.Error)
		return
	}
	appEnv.Render.JSON(w, http.StatusCreated, board)
}

// GetBoardHandler gets a board from the board store by id
func GetBoardHandler(w http.ResponseWriter, req *http.Request, appEnv AppEnv) {
	vars := mux.Vars(req)
	bid, _ := strconv.Atoi(vars["bid"])
	board := models.Board{}
	result := appEnv.DBConn.Preload("Sections").First(&board, bid)
	if result.Error != nil {
		appEnv.sendJSON(w, http.StatusNotFound, status.Response{
			Status:  strconv.Itoa(http.StatusNotFound),
			Message: "Can't find board",
		})
		return
	}
	appEnv.sendJSON(w, http.StatusOK, board)
}

// Update a board in the board store
func UpdateBoardHandler(w http.ResponseWriter, req *http.Request, appEnv AppEnv) {
	vars := mux.Vars(req)
	bid, _ := strconv.Atoi(vars["bid"])
	decoder := json.NewDecoder(req.Body)
	var b models.Board
	err := decoder.Decode(&b)
	if err != nil {
		appEnv.sendJSON(w, http.StatusBadRequest, status.Response{
			Status:  strconv.Itoa(http.StatusBadRequest),
			Message: "malformed board object",
		})
		return
	}

	board := models.Board{
		ID:          uint(bid),
		Title:       b.Title,
		Description: b.Description,
	}

	if err = appEnv.DBConn.Model(&board).Updates(&board).Error; err != nil {
		appEnv.sendInternalServerError(w, "Error updating board", err)
		return
	}

	appEnv.sendJSON(w, http.StatusOK, board)
}

// Delete a board from the board store
func DeleteBoardHandler(w http.ResponseWriter, req *http.Request, appEnv AppEnv) {
	vars := mux.Vars(req)
	bid, _ := strconv.Atoi(vars["bid"])
	board := models.Board{}
	result := appEnv.DBConn.Delete(&board, bid)
	if result.Error != nil {
		appEnv.sendInternalServerError(w, "Error deleting board", result.Error)
		return
	}
	// If the board was not found and due to it not being found it couldn't be deleted
	if result.RowsAffected == 0 {
		appEnv.sendJSON(w, http.StatusNotFound, status.Response{
			Status:  strconv.Itoa(http.StatusNotFound),
			Message: "Can't find board",
		})
		return
	}
	appEnv.sendJSON(w, http.StatusOK, board)
}
