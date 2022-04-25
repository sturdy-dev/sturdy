package main

import (
	"context"
	"log"

	"getsturdy.com/api/pkg/api"
	apiModule "getsturdy.com/api/pkg/api/module"
	"getsturdy.com/api/pkg/banner"
	xcontext "getsturdy.com/api/pkg/context"
	"getsturdy.com/api/pkg/datamigrations"
	"getsturdy.com/api/pkg/db"
	"getsturdy.com/api/pkg/di"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

func main() {
	app := func(c *di.Container) {
		c.Import(apiModule.Module)
		c.Import(xcontext.Module)
		c.Import(datamigrations.Module)
	}

	var (
		apiServer api.Starter
		ctx       context.Context
		dm        *datamigrations.Service
		sqlxdb    *sqlx.DB
		logger    *zap.Logger
	)

	if err := di.Init(app).To(&apiServer, &ctx, &dm, &sqlxdb, &logger); err != nil {
		log.Fatalf("failed to init: %+v", err)
	}

	if err := dm.Run(ctx); err != nil {
		logger.Fatal("failed to run data migrations", zap.Error(err))
	}

	if err := db.MigrateUP(sqlxdb.DB); err != nil {
		logger.Fatal("failed to migrate up", zap.Error(err))
	}

	banner.PrintBanner()

	if err := apiServer.Start(ctx); err != nil {
		logger.Fatal("failed to start server", zap.Error(err))
	}
}
