package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetColumnName(t *testing.T) {
	t.Run("should return column name", func(t *testing.T) {
		assert.Equal(t, "name", GetColumnName("Name"))
	})
	t.Run("should return column name", func(t *testing.T) {
		assert.Equal(t, "age", GetColumnName("Age"))
	})
	t.Run("shoudl return hard_name", func(t *testing.T) {
		assert.Equal(t, "hard_name", GetColumnName("HardName"))
	})
}
