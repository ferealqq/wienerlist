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

func CreateSectionHandler(baseController ctrl.BaseController[models.Section]) {
	var s models.Section
	if s, err := baseController.GetPostModel(s); err == nil {
		section := models.Section{
			Title:       s.Title,
			Description: s.Description,
			BoardId:     s.BoardId,
		}
		result := baseController.DB.Create(&section)
		if result.Error != nil {
			baseController.SendInternalServerError("Error creating a section", result.Error)
			return
		}
		baseController.SendJSON(http.StatusCreated, section)
	}
}
