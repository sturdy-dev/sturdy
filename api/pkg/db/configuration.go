package db

import (
	"time"

	"getsturdy.com/api/pkg/configuration/flags"
	"github.com/jmoiron/sqlx"
)

type Configuration struct {
	URL            flags.URL     `long:"url" description:"Database URL" required:"true" default:"postgres://mash:mash@127.0.0.1:5432/mash?sslmode=disable"`
	ConnectTimeout time.Duration `long:"connect-timeout" description:"Maximum time to wait for a connection to the database" default:"5s"`
}

func FromConfiguration(configuration *Configuration) (*sqlx.DB, error) {
	return SetupWithTimeout(configuration.URL.String(), configuration.ConnectTimeout)
}
