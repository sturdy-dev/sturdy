package db

import "getsturdy.com/api/pkg/di"

func Module(c *di.Container) {
	c.Register(NewGitHubInstallationRepo)
	c.Register(NewGitHubPRRepo)
	c.Register(NewGitHubRepositoryRepo)
	c.Register(NewGitHubUserRepo)
}
