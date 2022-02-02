//go:build enterprise
// +build enterprise

package module

import (
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/installations/db"
	"getsturdy.com/api/pkg/installations/enterprise/selfhosted"
	"getsturdy.com/api/pkg/installations/global"
	"getsturdy.com/api/pkg/installations/service"
)

func Module(c *di.Container) {
	c.Import(db.Module)
	c.Import(global.Module)
	c.Import(service.Module)
	c.Import(selfhosted.Module)
}
