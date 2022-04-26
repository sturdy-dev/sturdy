package waitinglist

import (
	"getsturdy.com/api/pkg/db/module"
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/waitinglist/acl"
	"getsturdy.com/api/pkg/waitinglist/instantintegration"
)

func Module(c *di.Container) {
	c.Import(module.Module)
	c.Register(NewWaitingListRepository)
	c.Register(acl.NewACLInterestRepository)
	c.Register(instantintegration.NewInstantIntegrationInterestRepository)
}
