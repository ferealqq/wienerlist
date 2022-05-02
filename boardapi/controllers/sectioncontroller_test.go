package controllers

import (
	"encoding/json"
	"net/http"
	"testing"

	. "github.com/ferealqq/golang-trello-copy/server/boardapi/models"
	app "github.com/ferealqq/golang-trello-copy/server/pkg/appenv"
	ctrl "github.com/ferealqq/golang-trello-copy/server/pkg/controller"
	"github.com/ferealqq/golang-trello-copy/server/pkg/database"
	. "github.com/ferealqq/golang-trello-copy/server/pkg/testing"
	. "github.com/ferealqq/golang-trello-copy/server/seeders"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestListSectionsHandler(t *testing.T) {
	action := HttpTestAction[Board]{
		Method: http.MethodGet,
		RouterFunc: func(e *gin.Engine, ae app.AppEnv) {
			e.GET("/sections", ctrl.MakeHandler(ae, ListSectionsHandler))
		},
		ReqPath: "/sections",
		Seeders: []func(db *gorm.DB){SeedSections},
		Tables:  []string{"sections"},
	}
	response := action.Run()
	assert.Equal(t, http.StatusOK, response.Code, "they should be equal")
	var d map[string]interface{}

	if err := json.Unmarshal(response.Body.Bytes(), &d); err != nil {
		assert.Fail(t, "Unmarshal should not fail")
		return
	}
	var count int64
	err := database.DBConn.Model(&Section{}).Count(&count)
	if err.Error != nil {
		assert.Fail(t, "Fetching count should not fail")
	}
	assert.Equal(t, int(count), int(d["count"].(float64)), "they should be equal")
}

/**
// FIXME Does not work with the current HttpTestAction wrapper. We need to refactor it to suit this functions needs.
func TestCreateSectionHandler(t *testing.T) {
	s := Section{
		Title:       "Test section",
		Description: "This is a test section",
	}
	b, _ := json.Marshal(s)

	action := HttpTestAction[Section]{
		Method: http.MethodPost,
		RouterFunc: func(e *gin.Engine, ae app.AppEnv) {
			e.POST("/sections", ctrl.MakeHandler(ae, CreateSectionHandler))
		},
		Body:    bytes.NewReader(b),
		Seeders: []func(db *gorm.DB){SeedSections},
		Tables:  []string{"sections"},
	}

	action.Run()

	var section Section
	err := database.DBConn.First(&section).Error
	if err != nil {
		assert.Fail(t, "Fetching section should not fail")
	}
	assert.Equal(t, section.Title, section.Title, "they should be equal")
	assert.Equal(t, section.Description, section.Description, "they should be equal")
	assert.Equal(t, section.BoardId, section.BoardId, "they should be equal")
}
*/
