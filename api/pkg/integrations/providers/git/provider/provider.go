package provider

import (
	"getsturdy.com/api/pkg/integrations/providers"
)

type Provider struct {
}

func New() *Provider {
	p := &Provider{}
	providers.Register(p)
	return p
}

func (p *Provider) ProviderType() providers.ProviderType {
	return providers.ProviderTypePushPull
}

func (p *Provider) ProviderName() providers.ProviderName {
	return providers.ProviderNameGit
}
