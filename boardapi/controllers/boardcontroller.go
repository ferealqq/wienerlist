package controllers

import (
	"net/http"
	"strconv"

	models "github.com/ferealqq/golang-trello-copy/server/boardapi/models"
	ctrl "github.com/ferealqq/golang-trello-copy/server/pkg/controller"
	"github.com/ferealqq/golang-trello-copy/server/pkg/health"
	"github.com/ferealqq/golang-trello-copy/server/pkg/status"
)

// HandlerFunc is a custom implementation of the http.HandlerFunc

// HealthcheckHandler returns useful info about the app
func HealthcheckHandler(baseController ctrl.BaseController[models.Board]) {
	check := health.Check{
		AppName: "golang-trello-copy",
		Version: baseController.AppEnv.Version,
	}
	baseController.SendJSON(http.StatusOK, check)
}

func ListBoardsHandler(baseController ctrl.BaseController[models.Board]) {
	var boards []models.Board
	result := baseController.DB.Preload("Sections").Find(&boards)
	if result.Error != nil {
		baseController.SendInternalServerError("Error listing boards", result.Error)
		return
	}
	responseObject := make(map[string]interface{})
	responseObject["bords"] = boards
	responseObject["count"] = len(boards)
	baseController.SendJSON(http.StatusOK, responseObject)
}

func CreateBoardHandler(baseController ctrl.BaseController[models.Board]) {
	// TODO Validation
	var b models.Board
	if b, err := baseController.GetPostModel(b); err == nil {
		board := models.Board{
			// FIXME: Couldn't we just give the b model to the model?
			Title:       b.Title,
			Description: b.Description,
		}
		result := baseController.DB.Create(&board)
		if result.Error != nil {
			baseController.SendInternalServerError("Error creating a board", result.Error)
			return
		}
		baseController.SendJSON(http.StatusCreated, board)
	}
}

// GetBoardHandler gets a board from the board store by id
func GetBoardHandler(baseController ctrl.BaseController[models.Board]) {
	if ID, err := baseController.GetUriId(); err == nil {
		board := models.Board{}
		result := baseController.DB.Preload("Sections").First(&board, ID)
		if result.Error != nil {
			baseController.SendJSON(http.StatusNotFound, status.Response{
				Status:  strconv.Itoa(http.StatusNotFound),
				Message: "Can't find board",
			})
			return
		}
		baseController.SendJSON(http.StatusOK, board)
	}
}

// Update a board in the board store
func UpdateBoardHandler(baseController ctrl.BaseController[models.Board]) {
	if bid, err := baseController.GetUriId(); err == nil {
		// TODO this should be a reusable function, used twice in this file
		var b models.Board
		if b, err := baseController.GetPostModel(b); err == nil {
			board := models.Board{
				ID:          uint(bid),
				Title:       b.Title,
				Description: b.Description,
			}

			if err = baseController.DB.Model(&board).Updates(&board).Error; err != nil {
				baseController.SendInternalServerError("Error updating board", err)
				return
			}

			baseController.SendJSON(http.StatusOK, board)
		}
	}
}

// Delete a board from the board store
func DeleteBoardHandler(baseController ctrl.BaseController[models.Board]) {
	if ID, err := baseController.GetUriId(); err == nil {
		var board models.Board
		result := baseController.DB.Delete(&board, ID)
		if result.Error != nil {
			baseController.SendInternalServerError("Error deleting board", result.Error)
			return
		}
		// If the board was not found and due to it not being found it couldn't be deleted
		if result.RowsAffected == 0 {
			baseController.SendJSON(http.StatusNotFound, status.Response{
				Status:  strconv.Itoa(http.StatusNotFound),
				Message: "Can't find board",
			})
			return
		}
		baseController.SendJSON(http.StatusOK, board)
	}
}
