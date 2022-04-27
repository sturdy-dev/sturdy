package waitinglist

import (
	"getsturdy.com/api/pkg/db"
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/waitinglist/acl"
	"getsturdy.com/api/pkg/waitinglist/instantintegration"
)

func Module(c *di.Container) {
	c.Import(db.Module)
	c.Register(NewWaitingListRepository)
	c.Register(acl.NewACLInterestRepository)
	c.Register(instantintegration.NewInstantIntegrationInterestRepository)
}
