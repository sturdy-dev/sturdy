package enterprise

import (
	"mash/pkg/di"
	"mash/pkg/github/enterprise/client"
	"mash/pkg/github/enterprise/db"
	"mash/pkg/github/enterprise/service"
	"mash/pkg/github/enterprise/workers"
)

var Module = di.NewModule(
	di.Needs(db.Module),
	di.Needs(workers.Module),
	di.Needs(service.Module),

	di.Provides(func() (client.ClientProvider, client.PersonalClientProvider) {
		return client.NewClient, client.NewPersonalClient
	}),

	// todo: hack to solve circular dependency
	di.Provides(func() *service.ImporterQueue {
		return new(service.ImporterQueue)
	}),
	di.Invoke(func(sq *service.ImporterQueue, wq workers.ImporterQueue) {
		*sq = wq
	}),

	// todo: hack to solve circular dependency
	di.Provides(func() *service.ClonerQueue {
		return new(service.ClonerQueue)
	}),
	di.Invoke(func(sq *service.ClonerQueue, wq *workers.ClonerQueue) {
		*sq = wq
	}),
)
