package db

import (
	"fmt"

	"getsturdy.com/api/pkg/configuration"
	module_configuration "getsturdy.com/api/pkg/configuration/module"
	configuration_db "getsturdy.com/api/pkg/db/configuration"
	"getsturdy.com/api/pkg/di"
	"github.com/jmoiron/sqlx"
)

func Module(c *di.Container) {
	c.Import(module_configuration.Module)
	c.Register(FromConfiguration)
}

func TestModule(c *di.Container) {
	c.Import(configuration.TestModule)
	c.Register(func(config *configuration_db.Configuration) (*sqlx.DB, error) {
		db, err := FromConfiguration(config)
		if err != nil {
			return nil, fmt.Errorf("could not connect to database: %w", err)
		}
		if err := MigrateUP(db.DB); err != nil {
			return nil, fmt.Errorf("could not migrate up: %w", err)
		}
		return db, err
	})
}
