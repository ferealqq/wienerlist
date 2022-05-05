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

	if wsIds, success := base.Context.GetQueryArray("WorkspaceId"); success {
		result.Where("workspace_id IN ?", wsIds)
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

func CreateItemHandler(base ctrl.BaseController[models.Item]) {
	var i models.Item
	if err := base.GetPostModel(&i); err == nil {
		item := models.Item{
			Title:       i.Title,
			Description: i.Description,
			SectionId:   i.SectionId,
			WorkspaceId: i.WorkspaceId,
		}
		result := base.DB.Create(&item)
		if result.Error != nil {
			base.SendInternalServerError("Error creating a item", result.Error)
			return
		}
		base.SendJSON(http.StatusCreated, item)
	}
}

func UpdateItemHandler(base ctrl.BaseController[models.Item]) {
	if uriId, err := base.GetUriId(); err == nil {
		var i models.Item
		if err := base.GetPostModel(&i); err == nil {
			result := base.DB.Model(&models.Item{}).Where("id = ?", uriId).Updates(i)
			if result.Error != nil {
				base.SendInternalServerError("Error updating a item", result.Error)
				return
			}
			// Empty ok
			base.SendJSON(http.StatusOK, nil)
		}
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
