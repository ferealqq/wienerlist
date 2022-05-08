package controllers

import (
	"errors"
	"net/http"

	models "github.com/ferealqq/golang-trello-copy/server/boardapi/models"
	ctrl "github.com/ferealqq/golang-trello-copy/server/pkg/controller"
	"gorm.io/gorm"
)

func ListSectionsHandler(base ctrl.BaseController[models.Section]) {
	var sections []models.Section
	result := base.DB.
		Preload("Items").
		Limit(base.DefaultQueryInt("limit", 100)).
		Offset(base.DefaultQueryInt("skip", 0))
	if boardIds, success := base.Context.GetQueryArray("board_id"); success {
		result.Where("board_id IN ?", boardIds)
	}
	result = result.Find(&sections)

	if result.Error != nil {
		base.SendInternalServerError("Error listing sections", result.Error)
		return
	}
	base.SendJSON(http.StatusOK, map[string]interface{}{
		"sections": sections,
		"count":    len(sections),
	})
}

func CreateSectionHandler(base ctrl.BaseController[models.Section]) {
	var s models.Section
	if err := base.GetPostModel(&s); err == nil {
		section := models.Section{
			Title:       s.Title,
			Description: s.Description,
			BoardId:     s.BoardId,
		}
		result := base.DB.Create(&section)
		if result.Error != nil {
			base.SendInternalServerError("Error creating a section", result.Error)
			return
		}
		base.SendJSON(http.StatusCreated, section)
	}
}

func GetSectionHandler(base ctrl.BaseController[models.Section]) {
	if uriId, err := base.GetUriId(); err == nil {
		var s models.Section
		if err := base.DB.First(&s, uriId).Error; err == nil {
			base.SendJSON(http.StatusOK, s)
		} else if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
			base.SendNotFound("Section not found")
		} else {
			base.SendInternalServerError("Error getting section", err)
		}
	}
}

func UpdateSectionHandler(base ctrl.BaseController[models.Section]) {
	if uriId, err := base.GetUriId(); err == nil {
		var s models.Section
		if err := base.GetPostModel(&s); err == nil {
			if err := base.DB.Model(&s).Where("id = ?", uriId).Updates(&s).Error; err == nil {
				base.SendJSON(http.StatusOK, nil)
			} else {
				base.SendInternalServerError("Error updating section", err)
			}
		}
	}
}

func DeleteSectionHandler(base ctrl.BaseController[models.Section]) {
	if uriId, err := base.GetUriId(); err == nil {
		var section models.Section
		result := base.DB.Delete(&section, uriId)
		if result.Error != nil {
			base.SendInternalServerError("Error getting section", result.Error)
			return
		}

		if result.RowsAffected == 0 {
			base.SendNotFound("Section not found")
			return
		}

		base.SendJSON(http.StatusOK, section)
	}
}
