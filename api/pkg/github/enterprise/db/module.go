package db

import (
	"getsturdy.com/api/pkg/db/module"
	"getsturdy.com/api/pkg/di"
)

func Module(c *di.Container) {
	c.Import(module.Module)
	c.Register(NewGitHubInstallationRepository)
	c.Register(NewGitHubPRRepository)
	c.Register(NewGitHubRepositoryRepository)
	c.Register(NewGitHubUserRepository)
}
