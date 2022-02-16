package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"getsturdy.com/api/pkg/api"
	module_api "getsturdy.com/api/pkg/api/module"
	"getsturdy.com/api/pkg/banner"
	module_configuration "getsturdy.com/api/pkg/configuration/module"
	"getsturdy.com/api/pkg/di"
	module_snapshots "getsturdy.com/api/pkg/snapshots/module"
)

func main() {
	ctx := context.Background()
	mainModule := func(c *di.Container) {
		c.Register(func() context.Context {
			return ctx
		})

		c.Import(module_configuration.Module)
		c.Import(module_api.Module)
		c.Import(module_snapshots.Module)
	}

	var apiServer api.Starter
	if err := di.Init(&apiServer, mainModule); err != nil {
		fmt.Printf("%+v\n", err)
		os.Exit(1)
	}

	banner.PrintBanner()

	if err := apiServer.Start(ctx); err != nil {
		log.Fatalf("faild to start server: %s", err)
	}
}
