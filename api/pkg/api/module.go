package api

import (
	workers_ci "getsturdy.com/api/pkg/ci/workers"
	"getsturdy.com/api/pkg/di"
	worker_gc "getsturdy.com/api/pkg/gc/worker"
	"getsturdy.com/api/pkg/gitserver"
	"getsturdy.com/api/pkg/http"
	"getsturdy.com/api/pkg/metrics"
	"getsturdy.com/api/pkg/pprof"
	worker_snapshots "getsturdy.com/api/pkg/snapshots/worker"
)

func Module(c *di.Container) {
	c.Import(http.Module)
	c.Import(worker_snapshots.Module)
	c.Import(workers_ci.Module)
	c.Import(worker_gc.Module)
	c.Import(gitserver.Module)
	c.Import(pprof.Module)
	c.Import(metrics.Module)
	c.Register(ProvideAPI)
}
