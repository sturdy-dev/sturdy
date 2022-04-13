package db

import (
	"getsturdy.com/api/pkg/db"
	"getsturdy.com/api/pkg/di"
)

func Module(c *di.Container) {
	c.Import(db.Module)
	c.Register(NewGitHubInstallationRepository)
	c.Register(NewGitHubPRRepository)
	c.Register(NewGitHubRepositoryRepository)
	c.Register(NewGitHubUserRepository)
}
