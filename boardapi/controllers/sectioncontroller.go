package controllers

import (
	"net/http"

	models "github.com/ferealqq/golang-trello-copy/server/boardapi/models"
	ctrl "github.com/ferealqq/golang-trello-copy/server/pkg/controller"
)

func ListSectionsHandler(baseController ctrl.BaseController[models.Section]) {
	var sections []models.Section
	result := baseController.DB.
		Limit(baseController.OptionalQueryNumber("limit", 100)).
		Offset(baseController.OptionalQueryNumber("skip", 0)).
		Find(&sections)

	if result.Error != nil {
		baseController.SendInternalServerError("Error listing sections", result.Error)
		return
	}
	baseController.SendJSON(http.StatusOK, map[string]interface{}{
		"sections": sections,
		"count":    len(sections),
	})
}
