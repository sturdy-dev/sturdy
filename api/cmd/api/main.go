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
	"getsturdy.com/api/pkg/di"
)

func main() {
	app := func(c *di.Container) {
		c.Import(apiModule.Module)
		c.Import(xcontext.Module)
	}

	var apiServer api.Starter
	var ctx context.Context
	if err := di.Init(app).To(&apiServer, &ctx); err != nil {
		fmt.Printf("%+v\n", err)
		os.Exit(1)
	}

	banner.PrintBanner()

	if err := apiServer.Start(ctx); err != nil {
		log.Fatalf("faild to start server: %s", err)
	}
}
