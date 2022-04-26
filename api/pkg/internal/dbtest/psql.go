package dbtest

import (
	"context"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"

	db_migrator "getsturdy.com/api/pkg/db/migrator"
	"getsturdy.com/api/pkg/internal/sturdytest"
)

func DB(t *testing.T) *sqlx.DB {
	d, err := db_migrator.Setup(
		sturdytest.PsqlDbSourceForTesting(),
		&nopMigrator{},
	)
	assert.NoError(t, err)
	return d
}

func MustGetDB() *sqlx.DB {
	d, err := db_migrator.Setup(
		sturdytest.PsqlDbSourceForTesting(),
		&nopMigrator{},
	)
	if err != nil {
		panic(err)
	}
	return d
}

type nopMigrator struct{}

func (n nopMigrator) Run(ctx context.Context, currentDatabaseSchemaVersion uint) error {
	return nil
}
