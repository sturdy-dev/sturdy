package main

import (
	"context"
	"log"

	"go.uber.org/zap"

	"getsturdy.com/api/pkg/api"
	apiModule "getsturdy.com/api/pkg/api/module"
	"getsturdy.com/api/pkg/banner"
	xcontext "getsturdy.com/api/pkg/context"
	"getsturdy.com/api/pkg/db/migrate"
	"getsturdy.com/api/pkg/di"
)

func main() {
	// run migrations
	var migrateService *migrate.Service
	if err := di.Init(migrate.Module).To(&migrateService); err != nil {
		log.Fatalf("failed to init: %+v", err)
	}

	if err := migrateService.Migrate(context.Background()); err != nil {
		log.Fatalf("failed to migrate up: %+v", err)
	}

	// start the application
	app := func(c *di.Container) {
		c.Import(apiModule.Module)
		c.Import(xcontext.Module)
	}

	var (
		ctx       context.Context
		apiServer api.Starter
		logger    *zap.Logger
	)
	if err := di.Init(app).To(&ctx, &apiServer, &logger); err != nil {
		log.Fatalf("failed to init: %+v", err)
	}

	banner.PrintBanner()

	if err := apiServer.Start(ctx); err != nil {
		logger.Fatal("failed to start server", zap.Error(err))
	}
}
