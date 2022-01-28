package waitinglist

import (
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/waitinglist/acl"
	"getsturdy.com/api/pkg/waitinglist/instantintegration"
)

func Module(c *di.Container) {
	c.Register(NewWaitingListRepo)
	c.Register(acl.NewACLInterestRepo)
	c.Register(instantintegration.NewInstantIntegrationInterestRepo)
}
