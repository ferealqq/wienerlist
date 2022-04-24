package boardapi

import (
	"encoding/json"
	"net/http"
	"strconv"

	models "github.com/ferealqq/golang-trello-copy/server/boardapi/models"
	"github.com/ferealqq/golang-trello-copy/server/pkg/health"
	"github.com/ferealqq/golang-trello-copy/server/pkg/status"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
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
	appEnv.Render.JSON(w, http.StatusOK, check)
}

func ListBoardsHandler(w http.ResponseWriter, req *http.Request, appEnv AppEnv) {
	var boards []models.Board
	result := appEnv.DBConn.Find(&boards)
	if result.Error != nil {
		response := status.Response{
			Status:  strconv.Itoa(http.StatusInternalServerError),
			Message: "Error fetching boards",
		}
		log.WithFields(log.Fields{
			"env":    appEnv.Env,
			"status": http.StatusInternalServerError,
		}).Error("Error fetching boards")
		appEnv.Render.JSON(w, http.StatusInternalServerError, response)
		return
	}
	responseObject := make(map[string]interface{})
	responseObject["boards"] = boards
	responseObject["count"] = len(boards)
	appEnv.Render.JSON(w, http.StatusOK, responseObject)
}

func CreateBoardHandler(w http.ResponseWriter, req *http.Request, appEnv AppEnv) {
	decoder := json.NewDecoder(req.Body)
	var b models.Board
	err := decoder.Decode(&b)
	if err != nil {
		response := status.Response{
			Status:  strconv.Itoa(http.StatusBadRequest),
			Message: "malformed board object",
		}
		log.WithFields(log.Fields{
			"env":    appEnv.Env,
			"status": http.StatusBadRequest,
		}).Error("malformed board object")
		appEnv.Render.JSON(w, http.StatusBadRequest, response)
		return
	}
	board := models.Board{
		Title: 		 b.Title,
		Description: b.Description,
	}
	result := appEnv.DBConn.Create(&board)
	if result.Error != nil {
		response := status.Response{
			Status:  strconv.Itoa(http.StatusInternalServerError),
			Message: "error creating board",
		}
		log.WithFields(log.Fields{
			"env":    appEnv.Env,
			"status": http.StatusInternalServerError,
			"error": result.Error,
		}).Error("error creating board")
		appEnv.Render.JSON(w, http.StatusCreated, response)
		return
	}

	appEnv.Render.JSON(w, http.StatusCreated, board)
}

// GetBoardHandler gets a board from the board store by id
func GetBoardHandler(w http.ResponseWriter, req *http.Request, appEnv AppEnv) {
	vars := mux.Vars(req)
	bid, _ := strconv.Atoi(vars["bid"])
	board := models.Board{}
	result := appEnv.DBConn.First(&board, bid)
	if result.Error != nil {
		response := status.Response{
			Status:  strconv.Itoa(http.StatusNotFound),
			Message: "can't find board",
		}
		log.WithFields(log.Fields{
			"env":    appEnv.Env,
			"status": http.StatusNotFound,
			"error":  result.Error,
		}).Error("Can't find board")
		appEnv.Render.JSON(w, http.StatusNotFound, response)
		return
	}
	appEnv.Render.JSON(w, http.StatusOK, board)
}

// Update a board in the board store
func UpdateBoardHandler(w http.ResponseWriter, req *http.Request, appEnv AppEnv) {
	vars := mux.Vars(req)
	bid, _ := strconv.Atoi(vars["bid"])
	decoder := json.NewDecoder(req.Body)
	var b models.Board
	err := decoder.Decode(&b)
	if err != nil {
		response := status.Response{
			Status:  strconv.Itoa(http.StatusBadRequest),
			Message: "malformed board object",
		}
		log.WithFields(log.Fields{
			"env":    appEnv.Env,
			"status": http.StatusBadRequest,
		}).Error("malformed board object")
		appEnv.Render.JSON(w, http.StatusBadRequest, response)
		return
	}

	board := models.Board{
		ID: 		 uint(bid),
		Title: 		 b.Title,
		Description: b.Description,
	}

	if err = appEnv.DBConn.Model(&board).Updates(&board).Error ; err != nil {
		response := status.Response{
			Status:  strconv.Itoa(http.StatusInternalServerError),
			Message: "error updating board",
		}
		log.WithFields(log.Fields{
			"env":    appEnv.Env,
			"status": http.StatusInternalServerError,
			"error":  err,
		}).Error("error updating board")
		appEnv.Render.JSON(w, http.StatusInternalServerError, response)
		return
	}

	appEnv.Render.JSON(w, http.StatusOK, board)
}

// Delete a board from the board store
func DeleteBoardHandler(w http.ResponseWriter, req *http.Request, appEnv AppEnv) {
	vars := mux.Vars(req)
	bid, _ := strconv.Atoi(vars["bid"])
	board := models.Board{}
	result := appEnv.DBConn.Delete(&board, bid)
	if result.Error != nil {
		response := status.Response{
			Status:  strconv.Itoa(http.StatusInternalServerError),
			Message: "error deleting board",
		}
		log.WithFields(log.Fields{
			"env":    appEnv.Env,
			"status": http.StatusInternalServerError,
			"error": result.Error,
		}).Error("Error deleting board")
		appEnv.Render.JSON(w, http.StatusInternalServerError, response)
		return
	}
	// If the board was not found and due to it not being found it couldn't be deleted
	if result.RowsAffected == 0 {
		response := status.Response{
			Status:  strconv.Itoa(http.StatusNotFound),
			Message: "can't find board",
		}
		log.WithFields(log.Fields{
			"env":    appEnv.Env,
			"status": http.StatusNotFound,
		}).Error("Can't find board")
		appEnv.Render.JSON(w, http.StatusNotFound, response)
		return
	}
	appEnv.Render.JSON(w, http.StatusOK, board)
}