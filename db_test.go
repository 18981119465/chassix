package chassis

import (
	"testing"

	"github.com/stretchr/testify/assert"
	_ "gorm.io/driver/postgres"
	_ "gorm.io/driver/sqlite"

	"c6x.io/chassis/config"
)

func TestDBs(t *testing.T) {
	//defer CloseAllDB()
	// given
	config.LoadFromEnvFile()
	dbCfg := config.Databases()
	assert.NotEmpty(t, dbCfg)
	// when
	dbs, _ := DBs()
	// then
	assert.NotNil(t, dbs[1])
	db, _ := dbs[1].DB()
	assert.Nil(t, db.Ping())
}
