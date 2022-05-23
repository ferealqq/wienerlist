package controllers

import (
	"errors"
	"net/http"

	models "github.com/ferealqq/wienerlist/boardapi/models"
	ctrl "github.com/ferealqq/wienerlist/pkg/controller"
	"gorm.io/gorm"
)

func ListWorkspacesHandler(base ctrl.BaseController[models.Workspace]) {
	var wspaces []models.Workspace
	result := base.DB.
		Preload("Boards").
		Limit(base.DefaultQueryInt("limit", 100)).
		Offset(base.DefaultQueryInt("skip", 0)).
		Find(&wspaces)

	if result.Error != nil {
		base.SendInternalServerError("Error listing boards", result.Error)
		return
	}
	base.SendJSON(http.StatusOK, map[string]interface{}{
		"workspaces": wspaces,
		"count":      len(wspaces),
		//TODO Return limit & skip maybe?
	})
}

func CreateWorkspaceHandler(base ctrl.BaseController[models.Workspace]) {
	// TODO Validation?
	var w models.Workspace
	if err := base.GetPostModel(&w); err == nil {
		workspace := models.Workspace{
			// FIXME: Couldn't we just give the w model to the model?
			Title:       w.Title,
			Description: w.Description,
		}
		result := base.DB.Create(&workspace)
		if result.Error != nil {
			base.SendInternalServerError("Error creating a workspace", result.Error)
			return
		}
		base.SendJSON(http.StatusCreated, workspace)
	}
}

func GetWorkspaceHandler(base ctrl.BaseController[models.Workspace]) {
	if ID, err := base.GetUriId(); err == nil {
		space := models.Workspace{}
		result := base.DB.Preload("Boards").First(&space, ID)
		if result.Error != nil && errors.Is(result.Error, gorm.ErrRecordNotFound) {
			base.SendNotFound("Workspace not found")
			return
		} else if result.Error != nil {
			base.SendInternalServerError("Error getting workspace", result.Error)
			return
		}
		base.SendJSON(http.StatusOK, space)
	}
}

func UpdateWorkspaceHandler(base ctrl.BaseController[models.Workspace]) {
	if bid, err := base.GetUriId(); err == nil {
		// TODO this should be a reusable function, used twice in this file
		var w models.Workspace
		if err := base.GetPostModel(&w); err == nil {
			if err = base.DB.Model(&models.Workspace{}).Where("id = ?", bid).Updates(&w).Error; err != nil {
				base.SendInternalServerError("Error updating workspace", err)
				return
			}
			// Empty ok
			base.SendJSON(http.StatusOK, nil)
		}
	}
}

func DeleteWorkspaceHandler(base ctrl.BaseController[models.Workspace]) {
	if ID, err := base.GetUriId(); err == nil {
		var space models.Workspace
		result := base.DB.Delete(&space, ID)
		if result.Error != nil {
			base.SendInternalServerError("Error deleting space", result.Error)
			return
		}
		if result.RowsAffected == 0 {
			base.SendNotFound("Workspace not found")
			return
		}
		base.SendJSON(http.StatusOK, space)
	}
}
