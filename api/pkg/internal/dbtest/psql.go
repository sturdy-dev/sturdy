package dbtest

import (
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"

	"getsturdy.com/api/pkg/db"
	"getsturdy.com/api/pkg/db/migrate/schema"
	"getsturdy.com/api/pkg/internal/sturdytest"
)

func DB(t *testing.T) *sqlx.DB {
	sqldb, err := getDB()
	assert.NoError(t, err)
	return sqldb
}

func MustGetDB() *sqlx.DB {
	db, err := getDB()
	if err != nil {
		panic(err)
	}
	return db
}

func getDB() (*sqlx.DB, error) {
	sqldb, err := db.SetupWithTimeout(sturdytest.PsqlDbSourceForTesting(), time.Second*5)
	if err != nil {
		return nil, err
	}
	s, err := schema.New(sqldb)
	if err != nil {
		return nil, err
	}
	return sqldb, s.Up()
}
