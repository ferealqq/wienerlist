package controllers

import (
	"net/http"
	"strconv"

	models "github.com/ferealqq/golang-trello-copy/server/boardapi/models"
	. "github.com/ferealqq/golang-trello-copy/server/pkg/container"
	"github.com/ferealqq/golang-trello-copy/server/pkg/health"
	"github.com/ferealqq/golang-trello-copy/server/pkg/status"
	"github.com/gin-gonic/gin"
)

// HandlerFunc is a custom implementation of the http.HandlerFunc
type HandlerFunc func(*gin.Context, AppContainer)

// HealthcheckHandler returns useful info about the app
func HealthcheckHandler(context *gin.Context, appContainer AppContainer) {
	check := health.Check{
		AppName: "golang-trello-copy",
		Version: appContainer.Version,
	}
	appContainer.SendJSON(context, http.StatusOK, check)
}

func ListBoardsHandler(context *gin.Context, appContainer AppContainer) {
	var boards []models.Board
	result := appContainer.DBConn.Preload("Sections").Find(&boards)
	if result.Error != nil {
		appContainer.SendInternalServerError(context, "Error listing boards", result.Error)
		return
	}
	responseObject := make(map[string]interface{})
	responseObject["bords"] = boards
	responseObject["count"] = len(boards)
	appContainer.SendJSON(context, http.StatusOK, responseObject)
}

func CreateBoardHandler(context *gin.Context, appContainer AppContainer) {
	// TODO Validation
	var b models.Board
	if err := context.ShouldBindJSON(&b); err != nil {
		// TODO Form a pattern in which we want to return error's
		appContainer.SendJSON(context, http.StatusBadRequest, status.Response{
			Status:  strconv.Itoa(http.StatusBadRequest),
			Message: "malformed board object",
		})
		return
	}
	board := models.Board{
		Title:       b.Title,
		Description: b.Description,
	}
	result := appContainer.DBConn.Create(&board)
	if result.Error != nil {
		appContainer.SendInternalServerError(context, "Error creating a board", result.Error)
		return
	}
	appContainer.SendJSON(context, http.StatusCreated, board)
}

// GetBoardHandler gets a board from the board store by id
func GetBoardHandler(context *gin.Context, appContainer AppContainer) {
	if ID, err := appContainer.GetUriId(context); err == nil {
		board := models.Board{}
		result := appContainer.DBConn.Preload("Sections").First(&board, ID)
		if result.Error != nil {
			appContainer.SendJSON(context, http.StatusNotFound, status.Response{
				Status:  strconv.Itoa(http.StatusNotFound),
				Message: "Can't find board",
			})
			return
		}
		appContainer.SendJSON(context, http.StatusOK, board)
	}
}

// Update a board in the board store
func UpdateBoardHandler(context *gin.Context, appContainer AppContainer) {
	if bid, err := appContainer.GetUriId(context); err == nil {
		// TODO this should be a reusable function, used twice in this file
		var b models.Board
		if err := context.ShouldBindJSON(&b); err != nil {
			// TODO Form a pattern in which we want to return error's
			appContainer.SendJSON(context, http.StatusBadRequest, status.Response{
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

		if err = appContainer.DBConn.Model(&board).Updates(&board).Error; err != nil {
			appContainer.SendInternalServerError(context, "Error updating board", err)
			return
		}

		appContainer.SendJSON(context, http.StatusOK, board)
	}
}

// Delete a board from the board store
func DeleteBoardHandler(context *gin.Context, appContainer AppContainer) {
	if ID, err := appContainer.GetUriId(context); err == nil {
		var board models.Board
		result := appContainer.DBConn.Delete(&board, ID)
		if result.Error != nil {
			appContainer.SendInternalServerError(context, "Error deleting board", result.Error)
			return
		}
		// If the board was not found and due to it not being found it couldn't be deleted
		if result.RowsAffected == 0 {
			appContainer.SendJSON(context, http.StatusNotFound, status.Response{
				Status:  strconv.Itoa(http.StatusNotFound),
				Message: "Can't find board",
			})
			return
		}
		appContainer.SendJSON(context, http.StatusOK, board)
	}
}
