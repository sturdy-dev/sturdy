package client

import "getsturdy.com/api/pkg/di"

func Module(c *di.Container) {
	c.Register(func() InstallationClientProvider { return NewInstallationClient })
	c.Register(func() PersonalClientProvider { return NewPersonalClient })
	c.Register(func() AppClientProvider { return NewAppClient })
}
