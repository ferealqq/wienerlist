package services

import (
	"testing"

	"github.com/ferealqq/wienerlist/front/store/model"
	"github.com/stretchr/testify/assert"
)

func TestGetSectionsWithParam(t *testing.T) {
	api := NewApi("http://localhost:4000/api/v1")
	var secs model.ListSections
	// TODO Create wrapper actions to get section data
	assert.Nil(t, api.Params("board_id", 1).Get("/sections/").BindModel(&secs))
}
