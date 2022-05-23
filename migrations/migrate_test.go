package migrations

import (
	"testing"

	"github.com/ferealqq/wienerlist/pkg/database"
	"github.com/stretchr/testify/assert"
)

func TestMigrate(t *testing.T) {
	database.TestDBInit()
	MigrateSeedAfterwards(database.DBConn)
	defer database.Close()
	// FIXME rewrite when seeder funcs throw errors instead of panics
	assert.True(t, true, "should not panic")
}
