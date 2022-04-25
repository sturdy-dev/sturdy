package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"getsturdy.com/api/pkg/api"
	apiModule "getsturdy.com/api/pkg/api/module"
	"getsturdy.com/api/pkg/banner"
	xcontext "getsturdy.com/api/pkg/context"
	"getsturdy.com/api/pkg/datamigrations"
	"getsturdy.com/api/pkg/db"
	"getsturdy.com/api/pkg/di"

	"github.com/jmoiron/sqlx"
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
	)
	if err := di.Init(app).To(&apiServer, &ctx, &dm, &sqlxdb); err != nil {
		fmt.Printf("%+v\n", err)
		os.Exit(1)
	}

	if err := dm.Run(ctx); err != nil {
		log.Fatalf("failed to run data migrations: %s", err)
	}

	if err := db.MigrateUP(sqlxdb.DB); err != nil {
		log.Fatalf("failed to migrate up: %s", err)
	}

	banner.PrintBanner()

	if err := apiServer.Start(ctx); err != nil {
		log.Fatalf("faild to start server: %s", err)
	}
}
