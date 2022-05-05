package controllers

import (
	"net/http"

	models "github.com/ferealqq/golang-trello-copy/server/boardapi/models"
	ctrl "github.com/ferealqq/golang-trello-copy/server/pkg/controller"
)

func ListItemsHandler(base ctrl.BaseController[models.Item]) {
	var items []models.Item
	result := base.DB.
		Limit(base.DefaultQueryInt("limit", 100)).
		Offset(base.DefaultQueryInt("skip", 0))

	if boardIds, success := base.Context.GetQueryArray("BoardId"); success {
		result.Where("board_id IN ?", boardIds)
	}

	if sectionIds, success := base.Context.GetQueryArray("SectionId"); success {
		result.Where("section_id IN ?", sectionIds)
	}

	result = result.Find(&items)

	if result.Error != nil {
		base.SendInternalServerError("Error listing items", result.Error)
		return
	}
	base.SendJSON(http.StatusOK, map[string]interface{}{
		"items": items,
		"count": len(items),
	})
}

func CreateItemHandler(base ctrl.BaseController[models.Section]) {
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

/*
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
			section := models.Section{
				ID:          uriId,
				Title:       s.Title,
				Description: s.Description,
				BoardId:     s.BoardId,
			}

			if err := base.DB.Model(&section).Updates(&section).Error; err == nil {
				base.SendJSON(http.StatusOK, section)
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
*/
