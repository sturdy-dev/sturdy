package dbtest

import (
	"time"

	"github.com/jmoiron/sqlx"

	"getsturdy.com/api/pkg/db"
	"getsturdy.com/api/pkg/db/migrate/schema"
	"getsturdy.com/api/pkg/internal/sturdytest"
)

func DB() (*sqlx.DB, error) {
	sqldb, err := db.SetupWithTimeout(sturdytest.PsqlDbSourceForTesting(), time.Millisecond)
	if err != nil {
		return nil, err
	}
	s, err := schema.New(sqldb)
	if err != nil {
		return nil, err
	}
	return sqldb, s.Up()
}
