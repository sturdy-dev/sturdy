package module

import (
	"getsturdy.com/api/pkg/configuration"
	module_configuration "getsturdy.com/api/pkg/configuration/module"
	"getsturdy.com/api/pkg/db"
	"getsturdy.com/api/pkg/di"
)

func Module(c *di.Container) {
	c.Import(module_configuration.Module)
	c.Register(db.FromConfiguration)
}

func TestModule(c *di.Container) {
	c.Import(configuration.TestModule)
	// c.Register(func(config *configuration_db.Configuration, dm *datamigrations.Service) (*sqlx.DB, error) {
	// 	db, err := db.FromConfiguration(config)
	// 	if err != nil {
	// 		return nil, fmt.Errorf("could not connect to database: %w", err)
	// 	}
	// 	if err := db.MigrateUP(db.DB, dm); err != nil {
	// 		return nil, fmt.Errorf("could not migrate up: %w", err)
	// 	}
	// 	return db, err
	// })
}
