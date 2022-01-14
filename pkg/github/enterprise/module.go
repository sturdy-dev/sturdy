package enterprise

import (
	"mash/pkg/di"
	"mash/pkg/github/enterprise/client"
	"mash/pkg/github/enterprise/db"
	"mash/pkg/github/enterprise/service"
	"mash/pkg/github/enterprise/workers"
)

func Module(c *di.Container) {
	c.Import(workers.Module)
	c.Import(db.Module)
	c.Register(service.New)
	c.Register(func() (client.ClientProvider, client.PersonalClientProvider) {
		return client.NewClient, client.NewPersonalClient
	})

	// todo: hack to solve circular import dependency
	iq := new(service.ImporterQueue)
	c.Register(func() *service.ImporterQueue {
		return iq
	})

	type importerHack struct{}
	c.Register(func(wq workers.ImporterQueue) importerHack {
		*iq = wq
		return struct{}{}
	})

	// todo: hack to solve circular import dependency
	cq := new(service.ClonerQueue)
	c.Register(func() *service.ClonerQueue {
		return cq
	})
	type clonerHack struct{}
	c.Register(func(wq *workers.ClonerQueue) clonerHack {
		*cq = wq
		return struct{}{}
	})
}
