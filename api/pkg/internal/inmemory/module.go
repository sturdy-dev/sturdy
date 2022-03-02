package inmemory

import "getsturdy.com/api/pkg/di"

func TestModule(c *di.Container) {
	c.Register(NewInMemoryAclRepo)
	c.Register(NewInMemoryCodebaseRepo)
	c.Register(NewInMemoryCodebaseUserRepo)
	c.Register(NewInMemoryGitHubInstallationRepository)
	c.Register(NewInMemoryGitHubRepositoryRepo)
	c.Register(NewInMemoryGitHubUserRepo)
	c.Register(NewInMemoryOrganizationMemberRepository)
	c.Register(NewInMemoryOrganizationRepo)
	c.Register(NewInMemorySnapshotRepo)
	c.Register(NewInMemoryViewRepo)
}
