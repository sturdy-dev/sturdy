package provider

import (
	db_acl "getsturdy.com/api/pkg/codebases/acl/db"
	db_codebases "getsturdy.com/api/pkg/codebases/db"
	"getsturdy.com/api/pkg/di"
	service_users "getsturdy.com/api/pkg/users/service/module"
)

func Module(c *di.Container) {
	c.Import(db_acl.Module)
	c.Import(db_codebases.Module)
	c.Import(service_users.Module)
	c.Register(New)
}
