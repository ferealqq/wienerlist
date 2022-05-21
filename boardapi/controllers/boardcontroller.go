package controllers

import (
	"errors"
	"net/http"

	models "github.com/ferealqq/wienerlist/boardapi/models"
	ctrl "github.com/ferealqq/wienerlist/pkg/controller"
	"github.com/ferealqq/wienerlist/pkg/health"
	"gorm.io/gorm"
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
	result := baseController.DB.
		Preload("Sections").
		Limit(baseController.DefaultQueryInt("limit", 100)).
		Offset(baseController.DefaultQueryInt("skip", 0)).
		Find(&boards)

	if result.Error != nil {
		baseController.SendInternalServerError("Error listing boards", result.Error)
		return
	}
	baseController.SendJSON(http.StatusOK, map[string]interface{}{
		"boards": boards,
		"count":  len(boards),
	})
}

func CreateBoardHandler(baseController ctrl.BaseController[models.Board]) {
	// TODO Validation?
	var b models.Board
	if err := baseController.GetPostModel(&b); err == nil {
		board := models.Board{
			// FIXME: Couldn't we just give the b model to the model?
			Title:       b.Title,
			Description: b.Description,
			WorkspaceId: b.WorkspaceId,
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
		if result.Error != nil && errors.Is(result.Error, gorm.ErrRecordNotFound) {
			baseController.SendNotFound("Board not found")
			return
		} else if result.Error != nil {
			baseController.SendInternalServerError("Error getting board", result.Error)
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
		if err := baseController.GetPostModel(&b); err == nil {
			if err = baseController.DB.Model(&models.Board{}).Where("id = ?", bid).Updates(&b).Error; err != nil {
				baseController.SendInternalServerError("Error updating board", err)
				return
			}
			// Empty ok
			baseController.SendJSON(http.StatusOK, nil)
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
			baseController.SendNotFound("Board not found")
			return
		}
		baseController.SendJSON(http.StatusOK, board)
	}
}
