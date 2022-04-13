package enterprise

import (
	db_codebases "getsturdy.com/api/pkg/codebases/db"
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/emails/transactional"
	db_github "getsturdy.com/api/pkg/github/enterprise/db"
)

func Module(c *di.Container) {
	c.Import(transactional.Module)
	c.Import(db_github.Module)
	c.Import(db_codebases.Module)
	c.Register(New, new(transactional.EmailSender))
}
