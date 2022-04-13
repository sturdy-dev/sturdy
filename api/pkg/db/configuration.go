package db

import (
	"getsturdy.com/api/pkg/db/configuration"

	"github.com/jmoiron/sqlx"
)

func FromConfiguration(configuration *configuration.Configuration) (*sqlx.DB, error) {
	return SetupWithTimeout(configuration.URL.String(), configuration.ConnectTimeout)
}
