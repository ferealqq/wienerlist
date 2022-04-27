package boardapi

import (
	"net/http"
	"strconv"

	models "github.com/ferealqq/golang-trello-copy/server/boardapi/models"
	"github.com/ferealqq/golang-trello-copy/server/pkg/health"
	"github.com/ferealqq/golang-trello-copy/server/pkg/status"
	"github.com/gin-gonic/gin"
)

// HandlerFunc is a custom implementation of the http.HandlerFunc
type HandlerFunc func(*gin.Context, AppEnv)

func MakeHandler(appEnv AppEnv, fn func(*gin.Context, AppEnv)) func(*gin.Context) {
	return func(c *gin.Context) {
		fn(c, appEnv)
	}
}

// HealthcheckHandler returns useful info about the app
func HealthcheckHandler(context *gin.Context, appEnv AppEnv) {
	check := health.Check{
		AppName: "golang-trello-copy",
		Version: appEnv.Version,
	}
	appEnv.sendJSON(context, http.StatusOK, check)
}

func ListBoardsHandler(context *gin.Context, appEnv AppEnv) {
	var boards []models.Board
	result := appEnv.DBConn.Preload("Sections").Find(&boards)
	if result.Error != nil {
		appEnv.sendInternalServerError(context, "Error listing boards", result.Error)
		return
	}
	responseObject := make(map[string]interface{})
	responseObject["bords"] = boards
	responseObject["count"] = len(boards)
	appEnv.sendJSON(context, http.StatusOK, responseObject)
}

func CreateBoardHandler(context *gin.Context, appEnv AppEnv) {
	// TODO Validation
	var b models.Board
	if err := context.ShouldBindJSON(&b); err != nil {
		// TODO Form a pattern in which we want to return error's
		appEnv.sendJSON(context, http.StatusBadRequest, status.Response{
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
		appEnv.sendInternalServerError(context, "Error creating a board", result.Error)
		return
	}
	appEnv.sendJSON(context, http.StatusCreated, board)
}

// GetBoardHandler gets a board from the board store by id
func GetBoardHandler(context *gin.Context, appEnv AppEnv) {
	if ID, err := appEnv.getUriId(context); err == nil {
		board := models.Board{}
		result := appEnv.DBConn.Preload("Sections").First(&board, ID)
		if result.Error != nil {
			appEnv.sendJSON(context, http.StatusNotFound, status.Response{
				Status:  strconv.Itoa(http.StatusNotFound),
				Message: "Can't find board",
			})
			return
		}
		appEnv.sendJSON(context, http.StatusOK, board)
	}
}

// Update a board in the board store
func UpdateBoardHandler(context *gin.Context, appEnv AppEnv) {
	if bid, err := appEnv.getUriId(context); err == nil {
		// TODO this should be a reusable function, used twice in this file
		var b models.Board
		if err := context.ShouldBindJSON(&b); err != nil {
			// TODO Form a pattern in which we want to return error's
			appEnv.sendJSON(context, http.StatusBadRequest, status.Response{
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
			appEnv.sendInternalServerError(context, "Error updating board", err)
			return
		}

		appEnv.sendJSON(context, http.StatusOK, board)
	}
}

// Delete a board from the board store
func DeleteBoardHandler(context *gin.Context, appEnv AppEnv) {
	if ID, err := appEnv.getUriId(context); err == nil {
		var board models.Board
		result := appEnv.DBConn.Delete(&board, ID)
		if result.Error != nil {
			appEnv.sendInternalServerError(context, "Error deleting board", result.Error)
			return
		}
		// If the board was not found and due to it not being found it couldn't be deleted
		if result.RowsAffected == 0 {
			appEnv.sendJSON(context, http.StatusNotFound, status.Response{
				Status:  strconv.Itoa(http.StatusNotFound),
				Message: "Can't find board",
			})
			return
		}
		appEnv.sendJSON(context, http.StatusOK, board)
	}
}
